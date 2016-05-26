package handlers

import (
	"encoding/json"
	"path"
	"os"
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"fmt"
)

const donePackingPartName string = "done-packing-part"
const packingPartName string = "packing-part"
func build(conn websocketConnection, data interface{}) {
	dat := data.(map[string]interface{})
	modpack := createSingleModpackData(dat["modpack"])
	mods := make([]ModResponse, 0)
	modsData := dat["mods"]
	modsDat, _ := json.Marshal(modsData)
	err := json.Unmarshal(modsDat, &mods)
	if err != nil {
		conn.Log("Could not create mod list")
		return
	}
	buildModpack(modpack, mods, conn)
}

func buildModpack(modpack Modpack, mods []ModResponse, conn websocketConnection) {
	// Create output directory
	outputDirectory := path.Join(modpack.OutputDirectory, modpack.Name)
	os.MkdirAll(outputDirectory, os.ModePerm)
	var total int

	// Handle forge
	if modpack.Technic.CreateForgeZip {
		total++
		go packForgeFolder(modpack, conn, outputDirectory)
	}

	// Handle any additional folder
	for _, folder := range modpack.AdditionalFolders {
		if folder.Include {
			total++
			go packAdditionalFolder(modpack, path.Join(modpack.InputDirectory, folder.Name), outputDirectory, conn)
		}
	}

	// Handle mods
	total += len(mods)
	Write(conn, "total-to-pack", total)
	for _, mod := range mods {
		go packMod(mod, conn, outputDirectory)
	}
}

func packForgeFolder(modpack Modpack, conn websocketConnection, outputDirectory string) {
	const minecraftForge string = "Minecraft Forge"
	Write(conn, packingPartName, minecraftForge)
	outputDirectory = path.Join(outputDirectory, "mods", "forge")
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, "forge-" + modpack.MinecraftVersion + "-" + fmt.Sprintf("%v", modpack.Technic.ForgeVersion.Build) + ".zip")

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
		return;
	}
	Write(conn, donePackingPartName, minecraftForge)
}

func packAdditionalFolder(modpack Modpack, folderPath string, outputDirectory string, conn websocketConnection) {
	Write(conn, packingPartName, folderPath)
	inputFolderInfo, _ := os.Stat(folderPath)
	outputDirectory = path.Join(outputDirectory, "mods", modpack.Name + "-" + inputFolderInfo.Name())
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, modpack.Name + "-" + inputFolderInfo.Name() + "-" + modpack.GetVersionString() + ".zip")
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
				return;
			}

			f, err := zipWriter.Create(zipName)
			if err != nil {
				conn.Log("Error while creating zip file content: " + err.Error() + "\n" + folder + "/" + file.Name())
				return
			}

			_, err = io.Copy(f, fileReader)
			if err != nil {
				conn.Log("Error white writing content to zip file: " + err.Error() + "\n" + folder + "/" + file.Name())
				return;
			}
		}
	}
}

func packMod(mod ModResponse, conn websocketConnection, outputDirectory string) {
	Write(conn, packingPartName, mod.Filename)
	outputDirectory = path.Join(outputDirectory, "mods", mod.ModId)
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, mod.getVersionString() + ".zip")
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
}
