package handlers

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/zlepper/go-modpack-packer/source/backend/solder"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/s3"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const donePackingPartName string = "done-packing-part"
const packingPartName string = "packing-part"

func build(conn types.WebsocketConnection, data interface{}) {
	dat := data.(map[string]interface{})
	modpack := types.CreateSingleModpackData(dat["modpack"])
	mods := make([]types.Mod, 0)
	modsData := dat["mods"]
	modsDat, _ := json.Marshal(modsData)
	err := json.Unmarshal(modsDat, &mods)
	if err != nil {
		conn.Log("Could not create mod list")
		return
	}
	buildModpack(modpack, mods, conn)
}

func buildModpack(modpack types.Modpack, mods []types.Mod, conn types.WebsocketConnection) {
	// Create output directory
	outputDirectory := path.Join(modpack.OutputDirectory, modpack.Name)
	os.MkdirAll(outputDirectory, os.ModePerm)
	var total int

	var ch chan *types.OutputInfo
	ch = make(chan *types.OutputInfo)

	// Handle forge
	if modpack.Technic.CreateForgeZip {
		total++
		go packForgeFolder(modpack, conn, outputDirectory, &ch)
	}

	// Handle any additional folder
	for _, folder := range modpack.AdditionalFolders {
		if folder.Include {
			total++
			go packAdditionalFolder(modpack, path.Join(modpack.InputDirectory, folder.Name), outputDirectory, conn, &ch)
		}
	}

	// Handle mods
	total += len(mods)
	conn.Write("total-to-pack", total)
	for _, mod := range mods {
		go packMod(mod, conn, outputDirectory, &ch)
	}

	infos := make([]*types.OutputInfo, 0)

	count := 0
	for count < total {
		info := <-ch
		infos = append(infos, info)
		count++
	}

	switch modpack.Technic.Upload.Type {
	case "none":
		{
			return // TODO Send back information and await clearance from client
		}
	case "ftp":
		{

		}
	case "s3":
		{
			s3.UploadFiles(&modpack, infos, conn)
		}
	}

	if modpack.Solder.Use {
		updateSolder(modpack, conn, infos)
	}
}

func ftpUpload(modpack types.Modpack) {

}

func updateSolder(modpack types.Modpack, conn types.WebsocketConnection, infos []*types.OutputInfo) {
	// Create solder client
	solderclient := solder.NewSolderClient(modpack.Solder.Url)
	loginSuccessful := solderclient.Login(modpack.Solder.Username, modpack.Solder.Password)
	if !loginSuccessful {
		log.Panic("Could not login to solder with the supplied credentials.")
	}

	var modpackId string
	if solderclient.IsPackOnline(&modpack) {
		modpackId = solderclient.GetModpackId(modpack.GetSlug())
	} else {
		modpackId = solderclient.CreatePack(modpack.Name, modpack.GetSlug())
	}

	var buildId string
	if solderclient.IsBuildOnline(&modpack) {
		buildId = solderclient.GetBuildId(&modpack)
	} else {
		buildId = solderclient.CreateBuild(&modpack, modpackId)
	}

	for _, info := range infos {
		go addInfoToSolder(info, buildId, conn, solderclient)
	}
}

func addInfoToSolder(info *types.OutputInfo, buildId string, conn types.WebsocketConnection, solderclient *solder.SolderClient) {
	var modid string
	modid = solderclient.GetModId(info.Id)
	if modid == "" {
		modid = solderclient.AddMod(info)
	}
	if modid == "" {
		log.Panic("Something went wrong wehn adding a mod to solder.")
	}

	if !solderclient.IsModversionOnline(info) {
		md5, err := ComputeMd5(info.File)
		if err != nil {
			log.Panic(err)
		}
		solderclient.AddModVersion(info.Id, hex.EncodeToString(md5), info.GenerateOnlineVersion())
	}
	if !solderclient.IsModversionActiveInBuild(info, buildId) {
		if solderclient.IsModInBuild(info, buildId) {
			solderclient.SetModVersionInBuild(info, buildId)
		} else {
			solderclient.AddModversionToBuild(info, buildId)
		}
	}

	conn.Write("done-updating-solder", info.ProgressKey)
}

func ComputeMd5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}

func packForgeFolder(modpack types.Modpack, conn types.WebsocketConnection, outputDirectory string, ch *chan *types.OutputInfo) {
	const minecraftForge string = "Minecraft Forge"
	outputDirectory = path.Join(outputDirectory, "mods", "forge")
	os.MkdirAll(outputDirectory, os.ModePerm)
	version := fmt.Sprintf("%v", modpack.Technic.ForgeVersion.Build)
	outputFile := path.Join(outputDirectory, "forge-"+modpack.MinecraftVersion+"-"+version+".zip")
	conn.Write(packingPartName, minecraftForge)

	zipfile, err := os.Create(outputFile)
	if err != nil {
		conn.Log("Error when creating zip file: " + err.Error() + "\n" + outputFile)
		return
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	resp, err := http.Get(modpack.Technic.ForgeVersion.DownloadUrl)
	if err != nil {
		conn.Log("Error downloading forge file: " + err.Error() + "\n" + modpack.Technic.ForgeVersion.DownloadUrl)
		return
	}
	defer resp.Body.Close()

	f, err := zipWriter.Create("bin/modpack.jar")
	if err != nil {
		conn.Log("Error while creating zip file content: " + err.Error())
		return
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		conn.Log("Error white writing content to zip file: " + err.Error())
		return
	}
	conn.Write(donePackingPartName, minecraftForge)

	if ch != nil {
		info := types.OutputInfo{
			File:             outputFile,
			Name:             minecraftForge,
			Id:               "forge",
			Version:          version,
			MinecraftVersion: modpack.MinecraftVersion,
			Description:      "The core of everything modded minecraft.",
			Url:              "http://www.minecraftforge.net/",
			Author:           "LexManos, cpw",
			ProgressKey:      minecraftForge,
		}
		*ch <- &info
	}
}

func packAdditionalFolder(modpack types.Modpack, folderPath string, outputDirectory string, conn types.WebsocketConnection, ch *chan *types.OutputInfo) {
	conn.Write(packingPartName, folderPath)
	inputFolderInfo, _ := os.Stat(folderPath)
	outputDirectory = path.Join(outputDirectory, "mods", modpack.Name+"-"+inputFolderInfo.Name())
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, modpack.Name+"-"+inputFolderInfo.Name()+"-"+modpack.GetVersionString()+".zip")
	zipfile, err := os.Create(outputFile)
	if err != nil {
		conn.Log("Error when creating zip file: " + err.Error() + "\n" + outputFile)
		return
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	packFolder(zipWriter, folderPath, ".", conn)

	conn.Write(donePackingPartName, folderPath)

	if ch != nil {
		info := types.OutputInfo{
			File:             outputFile,
			Name:             modpack.Name + "-" + inputFolderInfo.Name(),
			Id:               strings.ToLower(modpack.Name + "-" + inputFolderInfo.Name()),
			Version:          modpack.Version,
			MinecraftVersion: modpack.MinecraftVersion,
			ProgressKey:      folderPath,
		}
		*ch <- &info
	}
}

func packFolder(zipWriter *zip.Writer, folder string, parent string, conn types.WebsocketConnection) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			packFolder(zipWriter, path.Join(folder, file.Name()), path.Join(parent, file.Name()), conn)
		} else {
			zipName := path.Join(parent, file.Name())

			fileReader, err := os.Open(path.Join(folder, file.Name()))
			defer fileReader.Close()

			if err != nil {
				conn.Log(err.Error())
				return
			}

			f, err := zipWriter.Create(zipName)
			if err != nil {
				conn.Log("Error while creating zip file content: " + err.Error() + "\n" + folder + "/" + file.Name())
				return
			}

			_, err = io.Copy(f, fileReader)
			if err != nil {
				conn.Log("Error white writing content to zip file: " + err.Error() + "\n" + folder + "/" + file.Name())
				return
			}
		}
	}
}

func packMod(mod types.Mod, conn types.WebsocketConnection, outputDirectory string, ch *chan *types.OutputInfo) {
	conn.Write(packingPartName, mod.Filename)
	outputDirectory = path.Join(outputDirectory, "mods", mod.ModId)
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, mod.GetVersionString()+".zip")
	zipfile, err := os.Create(outputFile)
	if err != nil {
		conn.Log("Error when creating zip file: " + err.Error() + "\n" + outputFile)
		return
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	fileInfo, _ := os.Stat(mod.Filename)
	file, _ := os.Open(mod.Filename)
	defer file.Close()

	zipName := path.Join("mods", fileInfo.Name())

	f, err := zipWriter.Create(zipName)
	if err != nil {
		conn.Log("Error while creating zip file content: " + err.Error() + "\n" + outputFile)
		return
	}

	_, err = io.Copy(f, file)
	if err != nil {
		conn.Log("Error white writing content to zip file: " + err.Error() + "\n" + outputFile)
	}

	conn.Write(donePackingPartName, mod.Filename)

	if ch != nil {
		info := types.OutputInfo{
			File:             outputFile,
			Name:             mod.Name,
			Id:               strings.ToLower(mod.ModId),
			Version:          mod.Version,
			MinecraftVersion: mod.MinecraftVersion,
			Description:      mod.Description,
			Author:           mod.Authors,
			Url:              mod.Url,
			ProgressKey:      mod.Filename,
		}
		*ch <- &info
	}
}
