package handlers

import (
	"archive/zip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/zlepper/go-modpack-packer/source/backend/db"
	"github.com/zlepper/go-modpack-packer/source/backend/helpers"
	"github.com/zlepper/go-modpack-packer/source/backend/solder"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/upload"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"github.com/zlepper/go-websocket-connection"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const donePackingPartName string = "done-packing-part"
const packingPartName string = "packing-part"

func build(conn websocket.WebsocketConnection, data interface{}) {
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

type uploadWaiting struct {
	Modpack types.Modpack       `json:"modpack"`
	Infos   []*types.OutputInfo `json:"infos"`
}

func continueRunning(conn websocket.WebsocketConnection, data interface{}) {
	dict := data.(map[string]interface{})

	var uploadInfo uploadWaiting
	err := mapstructure.Decode(dict, &uploadInfo)
	if err != nil {
		log.Panic(err)
	}

	solderclient, buildId := updateSolder(uploadInfo.Modpack)
	for _, info := range uploadInfo.Infos {
		go addInfoToSolder(info, buildId, conn, solderclient)
	}
}

func buildModpack(modpack types.Modpack, mods []types.Mod, conn websocket.WebsocketConnection) {
	// Create output directory
	outputDirectory := path.Join(modpack.OutputDirectory, modpack.Name)
	os.MkdirAll(outputDirectory, os.ModePerm)
	var total int

	var ch chan *types.OutputInfo
	ch = make(chan *types.OutputInfo)

	startTime := time.Now()
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
		// If the mod already is on solder, then we should likely skip it
		// however the user can override this. If they do we should still pack all files
		if !modpack.Technic.RepackAllMods && mod.IsOnSolder {
			continue
		}
		go packMod(mod, conn, outputDirectory, &ch)
	}

	infos := make([]*types.OutputInfo, 0)

	// Save the mods to the database
	d := db.GetModsDb()
	for i, _ := range mods {
		d.AddMod(&mods[i])
	}
	d.Save()

	count := 0
	for count < total {
		info := <-ch
		infos = append(infos, info)
		count++
	}
	endTime := time.Now()
	spendtime := endTime.UnixNano() - startTime.UnixNano()
	spendtime = spendtime / (int64(time.Millisecond) / int64(time.Nanosecond))
	fmt.Printf("Time spend packing: %d ms", spendtime)

	switch modpack.Technic.Upload.Type {
	case "none":
		{
			// Send back information and await clearance from client
			conn.Write("waiting-for-file-upload", uploadWaiting{Modpack: modpack, Infos: infos})
			return
		}
	case "ftp":
		{
			upload.UploadFilesToFtp(&modpack, infos, conn)
		}
	case "s3":
		{
			upload.UploadFilesToS3(&modpack, infos, conn)
		}
	}

	if modpack.Solder.Use {
		solderclient, buildId := updateSolder(modpack)

		for _, info := range infos {
			go addInfoToSolder(info, buildId, conn, solderclient)
		}
	}

}

func updateSolder(modpack types.Modpack) (*solder.SolderClient, string) {
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
	return solderclient, buildId
}

func addInfoToSolder(info *types.OutputInfo, buildId string, conn websocket.WebsocketConnection, solderclient *solder.SolderClient) {
	conn.Write("updating-solder", info.ProgressKey)
	var modid string
	modid = solderclient.GetModId(info.Id)
	if modid == "" {
		modid = solderclient.AddMod(info)
	}
	if modid == "" {
		log.Printf("%v\n", *info)
		log.Panic("Something went wrong wehn adding a mod to solder.")
	}

	md5, err := helpers.ComputeMd5(info.File)
	if !solderclient.IsModversionOnline(info) {
		if err != nil {
			log.Panic(err)
		}
		solderclient.AddModVersion(modid, hex.EncodeToString(md5), info.GenerateOnlineVersion())
	} else {
		id := solderclient.GetModVersionId(info)
		solderclient.RehashModVersion(id, hex.EncodeToString(md5))
	}
	if !solderclient.IsModversionActiveInBuild(info, buildId) {
		if solderclient.IsModInBuild(info, buildId) {
			solderclient.SetModVersionInBuild(info, buildId)
		} else {
			solderclient.AddModversionToBuild(info, buildId)
		}
	}
	go db.GetModsDb().MarkModAsOnSolder(hex.EncodeToString(md5))
	conn.Write("done-updating-solder", info.ProgressKey)
}

func packForgeFolder(modpack types.Modpack, conn websocket.WebsocketConnection, outputDirectory string, ch *chan *types.OutputInfo) {
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

func packAdditionalFolder(modpack types.Modpack, folderPath string, outputDirectory string, conn websocket.WebsocketConnection, ch *chan *types.OutputInfo) {
	conn.Write(packingPartName, folderPath)
	inputFolderInfo, _ := os.Stat(folderPath)
	s := safeNormalizeString(modpack.Name + "-" + inputFolderInfo.Name())
	outputDirectory = path.Join(outputDirectory, "mods", s)
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, s+"-"+modpack.GetVersionString()+".zip")
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
			Id:               s,
			Version:          modpack.Version,
			MinecraftVersion: modpack.MinecraftVersion,
			ProgressKey:      folderPath,
		}
		*ch <- &info
	}
}

func packFolder(zipWriter *zip.Writer, folder string, parent string, conn websocket.WebsocketConnection) {
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

func packMod(mod types.Mod, conn websocket.WebsocketConnection, outputDirectory string, ch *chan *types.OutputInfo) {
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
			Id:               safeNormalizeString(mod.ModId),
			Version:          mod.Version,
			MinecraftVersion: mod.MinecraftVersion,
			Description:      mod.Description,
			Author:           mod.Authors,
			ProgressKey:      mod.Filename,
			IsOnSolder:       mod.IsOnSolder,
		}
		u := mod.Url
		if len(u) > 0 && strings.Index(u, "http") != 0 {
			u = "http://" + u
		}
		info.Url = u
		*ch <- &info
	}
}

func safeNormalizeString(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "-", -1)
}
