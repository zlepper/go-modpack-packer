package handlers

import (
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"log"
	"os"
)

type inputDirData struct {
	InputDir string `json:"inputDir"`
}

func createInputDirData(data map[string]interface{}) inputDirData {
	var res inputDirData
	res.InputDir = data["inputDir"].(string)
	return res
}

func findAdditionalFolders(conn websocketConnection, data interface{}) {
	dir := createInputDirData(data.(map[string]interface{}))
	files, _ := ioutil.ReadDir(dir.InputDir)
	folders := []string{}
	// Iterate the files
	for _, file := range files {
		// We only need the directories
		if file.IsDir() {
			// The mods folder should be handles in a special way
			if file.Name() == "mods" {
				subFiles, _ := ioutil.ReadDir(filepath.Join(dir.InputDir, file.Name()))
				for _, subfile := range subFiles {
					if subfile.IsDir() {
						folders = append(folders, filepath.Join(file.Name(), subfile.Name()))
					}
				}
			} else {
				folders = append(folders, file.Name())
			}
		}
	}
	Write(conn, "found-folders", folders)
}

type TechnicConfig struct {
	IsSolderPack     bool `json:"isSolderPack"`
	CreateForgeZip   bool `json:"createForgeZip"`
	ForgeVersion     string `json:"forgeVersion"`
	CheckPermissions bool `json:"checkPermissions"`
	IsPublicPack     bool `json:"isPublicPack"`
}

type FtbConfig struct {
	IsPublicPack bool `json:"isPublicPack"`
}

type Modpack struct {
	Name                 string `json:"name"`
	InputDirectory       string `json:"inputDirectory"`
	OutputDirectory      string `json:"outputDirectory"`
	ClearOutputDirectory bool `json:"clearOutputDirectory"`
	MinecraftVersion     string `json:"minecraftVersion"`
	Version              string `json:"version"`
	AdditionalFolders    map[string]bool `json:"additionalFolders"`
	Technic              TechnicConfig `json:"technic"`
	Ftb                  FtbConfig `json:"ftb"`
}

func createModpackData(data interface{}) []Modpack {
	dL := data.([]interface{})
	modpacks := make([]Modpack, 0)
	for _, di := range dL {
		d := di.(map[string]interface{})
		//d := data.(map[string]interface{})
		modpack := Modpack{
			Name:d["name"].(string),
			InputDirectory:d["inputDirectory"].(string),
			OutputDirectory:d["outputDirectory"].(string),
			ClearOutputDirectory:d["clearOutputDirectory"].(bool),
			MinecraftVersion:d["minecraftVersion"].(string),
			Version:d["version"].(string),
			AdditionalFolders:make(map[string]bool, 0),
		}
		additionalFolders := d["additionalFolders"].(map[string]interface{})
		for folder, include := range additionalFolders {
			modpack.AdditionalFolders[folder] = include.(bool)
		}
		tConfigMap := d["technic"].(map[string]interface{})
		tConfig := TechnicConfig{
			IsSolderPack: tConfigMap["isSolderPack"].(bool),
			CreateForgeZip: tConfigMap["createForgeZip"].(bool),
			ForgeVersion:tConfigMap["forgeVersion"].(string),
			CheckPermissions:tConfigMap["checkPermissions"].(bool),
			IsPublicPack:tConfigMap["isPublicPack"].(bool),
		}
		modpack.Technic = tConfig
		modpack.Ftb = FtbConfig{
			IsPublicPack: d["ftb"].(map[string]interface{})["isPublicPack"].(bool),
		}
		modpacks = append(modpacks, modpack)
	}
	return modpacks
}

func saveModpacks(conn websocketConnection, data interface{}) {
	modpacks := createModpackData(data)
	modpackData, err := json.Marshal(modpacks)
	if err != nil {
		log.Panic(err)
	}
	// Get the appData directory, since go doesn't expose it, electron passes it as a parameter
	dataDirectory := os.Args[1]
	modpackFile := filepath.Join(dataDirectory, "modpacks.json")
	err = ioutil.WriteFile(modpackFile, modpackData, os.FileMode(0777))
	if err != nil {
		log.Panic(err)
	}
	log.Println("Data saved")

}

func loadModpacks(conn websocketConnection) {
	dataDirectory := os.Args[1]
	modpackFile := filepath.Join(dataDirectory, "modpacks.json")
	modpackData, err := ioutil.ReadFile(modpackFile)
	if err != nil {
		conn.Log("Unable to reload data " + err.Error())
		return
	}
	var modpacks []Modpack
	err = json.Unmarshal(modpackData, &modpacks)
	if err != nil {
		conn.Log("Could not parse json data " + err.Error())
		return
	}
	log.Printf("%v", modpacks)
	Write(conn, "data-loaded", modpacks)
}
