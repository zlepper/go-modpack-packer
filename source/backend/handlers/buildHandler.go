package handlers

import (
	"archive/zip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

const donePackingPartName string = "done-packing-part"
const packingPartName string = "packing-part"

func build(conn websocketConnection, data interface{}) {
	dat := data.(map[string]interface{})
	modpack := createSingleModpackData(dat["modpack"])
	mods := make([]Mod, 0)
	modsData := dat["mods"]
	modsDat, _ := json.Marshal(modsData)
	err := json.Unmarshal(modsDat, &mods)
	if err != nil {
		conn.Log("Could not create mod list")
		return
	}
	buildModpack(modpack, mods, conn)
}

type outputInfo struct {
	File             string
	Name             string
	Id               string
	Version          string
	MinecraftVersion string
	Description      string
	Author           string
	Url              string
	ProgressKey      string
}

func (o *outputInfo) GenerateOnlineVersion() string {
	return o.MinecraftVersion + "-" + o.Version
}

func buildModpack(modpack Modpack, mods []Mod, conn websocketConnection) {
	// Create output directory
	outputDirectory := path.Join(modpack.OutputDirectory, modpack.Name)
	os.MkdirAll(outputDirectory, os.ModePerm)
	var total int

	var ch chan outputInfo
	if modpack.Solder.Url {
		ch = make(chan outputInfo)
	}

	// Handle forge
	if modpack.Technic.CreateForgeZip {
		total++
		go packForgeFolder(modpack, conn, outputDirectory, ch)
	}

	// Handle any additional folder
	for _, folder := range modpack.AdditionalFolders {
		if folder.Include {
			total++
			go packAdditionalFolder(modpack, path.Join(modpack.InputDirectory, folder.Name), outputDirectory, conn, ch)
		}
	}

	// Handle mods
	total += len(mods)
	Write(conn, "total-to-pack", total)
	for _, mod := range mods {
		go packMod(mod, conn, outputDirectory, ch)
	}

	if ch != nil {
		updateSolder(modpack, ch, conn, total)
	}
}

func updateSolder(modpack Modpack, ch chan outputInfo, conn websocketConnection, total int) {
	// Create solder client
	solderclient := NewSolderClient(modpack.Solder.Url)
	loginSuccessful := solderclient.Login(modpack.Solder.Username, modpack.Solder.Password)
	if !loginSuccessful {
		panic("Could not login to solder with the supplied credentials.")
	}

	var modpackId string
	if solderclient.IsPackOnline(modpack) {
		modpackId = solderclient.GetModpackId(modpack.GetSlug())
	} else {
		modpackId = solderclient.CreatePack(modpack.Name, modpack.GetSlug())
	}

	var buildId string
	if solderclient.IsBuildOnline(modpack) {
		buildId = solderclient.GetBuildId(modpack)
	} else {
		buildId = solderclient.CreateBuild(modpack, modpackId)
	}
	count := 0
	for count < total {
		info := <-ch
		go addInfoToSolder(info, buildId, conn, solderclient)
		count++
	}

}

func addInfoToSolder(info outputInfo, buildId string, conn websocketConnection, solderclient *SolderClient) {
	var modid string
	modid = solderclient.GetModId(info.Id)
	if modid == "" {
		modid = solderclient.AddMod(info)
	}
	if modid == "" {
		panic("Something went wrong wehn adding a mod to solder.")
	}

	if !solderclient.IsModversionOnline(info) {
		md5, err := ComputeMd5(info.File)
		if err != nil {
			panic(err)
		}
		solderclient.AddModVersion(info.Id, string(md5), info.GenerateOnlineVersion())
	}
	if !solderclient.IsModversionActiveInBuild(info, buildId) {
		if solderclient.IsModInBuild(info, buildId) {
			solderclient.SetModVersionInBuild(info, buildId)
		} else {
			solderclient.AddModversionToBuild(info, buildId)
		}
	}

	Write(conn, "done-updating-solder", info.ProgressKey)
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

func packForgeFolder(modpack Modpack, conn websocketConnection, outputDirectory string, ch chan outputInfo) {
	const minecraftForge string = "Minecraft Forge"
	outputDirectory = path.Join(outputDirectory, "mods", "forge")
	os.MkdirAll(outputDirectory, os.ModePerm)
	version := fmt.Sprintf("%v", modpack.Technic.ForgeVersion.Build)
	outputFile := path.Join(outputDirectory, "forge-"+modpack.MinecraftVersion+"-"+version+".zip")
	Write(conn, packingPartName, minecraftForge)

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
	Write(conn, donePackingPartName, minecraftForge)

	if ch != nil {
		info := outputInfo{
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
		ch <- info
	}
}

func packAdditionalFolder(modpack Modpack, folderPath string, outputDirectory string, conn websocketConnection, ch chan outputInfo) {
	Write(conn, packingPartName, folderPath)
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

	Write(conn, donePackingPartName, folderPath)

	if ch != nil {
		info := outputInfo{
			File:             outputFile,
			Name:             modpack + "-" + inputFolderInfo.Name(),
			Id:               strings.ToLower(modpack + "-" + inputFolderInfo.Name()),
			Version:          modpack.Version,
			MinecraftVersion: modpack.MinecraftVersion,
			ProgressKey:      folderPath,
		}
		ch <- info
	}
}

func packFolder(zipWriter *zip.Writer, folder string, parent string, conn websocketConnection) {
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

func packMod(mod Mod, conn websocketConnection, outputDirectory string, ch chan outputInfo) {
	Write(conn, packingPartName, mod.Filename)
	outputDirectory = path.Join(outputDirectory, "mods", mod.ModId)
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, mod.getVersionString()+".zip")
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

	Write(conn, donePackingPartName, mod.Filename)

	if ch != nil {
		info := outputInfo{
			File:             outputFile,
			Name:             mod.Name,
			Id:               mod.ModId,
			Version:          mod.Version,
			MinecraftVersion: mod.MinecraftVersion,
			Description:      mod.Description,
			Author:           mod.Authors,
			Url:              mod.Url,
			ProgressKey:      mod.Filename,
		}
		ch <- info
	}
}
