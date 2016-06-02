package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

type ForgeVersion struct {
	Build            float64 `json:"build"`
	DownloadUrl      string  `json:"downloadUrl"`
	MinecraftVersion string  `json:"minecraftVersion"`
}

type TechnicConfig struct {
	IsSolderPack     float64      `json:"isSolderPack"`
	CreateForgeZip   bool         `json:"createForgeZip"`
	ForgeVersion     ForgeVersion `json:"forgeVersion"`
	CheckPermissions bool         `json:"checkPermissions"`
	IsPublicPack     bool         `json:"isPublicPack"`
}

type FtbConfig struct {
	IsPublicPack bool `json:"isPublicPack"`
}

type Folder struct {
	Name    string `json:"name"`
	Include bool   `json:"include"`
}

type SolderInfo struct {
	Use      bool   `json:"use"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Modpack struct {
	Name                 string        `json:"name"`
	InputDirectory       string        `json:"inputDirectory"`
	OutputDirectory      string        `json:"outputDirectory"`
	ClearOutputDirectory bool          `json:"clearOutputDirectory"`
	MinecraftVersion     string        `json:"minecraftVersion"`
	Version              string        `json:"version"`
	AdditionalFolders    []Folder      `json:"additionalFolders"`
	Technic              TechnicConfig `json:"technic"`
	Ftb                  FtbConfig     `json:"ftb"`
	Solder               SolderInfo    `json:"solder"`
	Memory               float64       `json:"memory"`
	Java                 string        `json:"java"`
}

func (m *Modpack) GetSlug() string {
	re := regexp.MustCompile("\\|/|\\||:|\\*|\"|<|>|\\?|'")
	s := re.ReplaceAllString(m.Name, "")
	s = strings.Replace(s, " ", "-", -1)
	s = strings.ToLower(s)
	return s
}

func (m *Modpack) GetVersionString() string {
	return m.MinecraftVersion + "-" + m.Version
}

func createSingleModpackData(di interface{}) Modpack {
	d := di.(map[string]interface{})
	//d := data.(map[string]interface{})
	modpack := Modpack{
		Name:                 d["name"].(string),
		InputDirectory:       d["inputDirectory"].(string),
		OutputDirectory:      d["outputDirectory"].(string),
		ClearOutputDirectory: d["clearOutputDirectory"].(bool),
		MinecraftVersion:     d["minecraftVersion"].(string),
		Version:              d["version"].(string),
		AdditionalFolders:    make([]Folder, 0),
	}
	additionalFolders := d["additionalFolders"].([]interface{})
	for _, folder := range additionalFolders {
		folderMap := folder.(map[string]interface{})

		f := Folder{
			Name:    folderMap["name"].(string),
			Include: folderMap["include"].(bool),
		}

		modpack.AdditionalFolders = append(modpack.AdditionalFolders, f)
	}
	tConfigMap := d["technic"].(map[string]interface{})
	tConfig := TechnicConfig{
		IsSolderPack:     tConfigMap["isSolderPack"].(float64),
		CreateForgeZip:   tConfigMap["createForgeZip"].(bool),
		CheckPermissions: tConfigMap["checkPermissions"].(bool),
		IsPublicPack:     tConfigMap["isPublicPack"].(bool),
	}
	fvInterface := tConfigMap["forgeVersion"]
	if fvInterface != nil {
		fvMap := tConfigMap["forgeVersion"].(map[string]interface{})
		fv := ForgeVersion{
			Build:            fvMap["build"].(float64),
			DownloadUrl:      fvMap["downloadUrl"].(string),
			MinecraftVersion: fvMap["minecraftVersion"].(string),
		}
		tConfig.ForgeVersion = fv
	}
	modpack.Technic = tConfig
	modpack.Ftb = FtbConfig{
		IsPublicPack: d["ftb"].(map[string]interface{})["isPublicPack"].(bool),
	}
	return modpack
}

func createModpackData(data interface{}) []Modpack {
	dL := data.([]interface{})
	modpacks := make([]Modpack, 0)
	for _, di := range dL {
		modpack := createSingleModpackData(di)
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
		conn.Log("Could not parse json data " + err.Error() + "\n" + string(modpackData))
		return
	}
	Write(conn, "data-loaded", modpacks)
}
