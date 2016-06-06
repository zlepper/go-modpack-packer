package handlers

import (
	"encoding/json"
	"github.com/zlepper/go-modpack-packer/source/backend/encryption"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type inputDirData struct {
	InputDir string `json:"inputDir"`
}

func createInputDirData(data map[string]interface{}) inputDirData {
	var res inputDirData
	res.InputDir = data["inputDir"].(string)
	return res
}

func findAdditionalFolders(conn types.WebsocketConnection, data interface{}) {
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
	conn.Write("found-folders", folders)
}

func saveModpacks(conn types.WebsocketConnection, data interface{}) {
	modpacks := types.CreateModpackData(data)
	for i, _ := range modpacks {
		modpack := &modpacks[i]
		modpack.Technic.Upload.AWS.SecretKey = encryption.EncryptString(modpack.Technic.Upload.AWS.SecretKey)
		modpack.Technic.Upload.AWS.AccessKey = encryption.EncryptString(modpack.Technic.Upload.AWS.AccessKey)
		modpack.Solder.Password = encryption.EncryptString(modpack.Solder.Password)
	}
	modpackData, err := json.Marshal(modpacks)
	if err != nil {
		log.Panic(err)
	}
	// Get the appData directory, since go doesn't expose it, electron passes it as a parameter
	dataDirectory := os.Args[1]
	modpackFile := filepath.Join(dataDirectory, "modpacks.json")
	err = ioutil.WriteFile(modpackFile, modpackData, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}

}

func loadModpacks(conn types.WebsocketConnection) {
	dataDirectory := os.Args[1]
	modpackFile := filepath.Join(dataDirectory, "modpacks.json")
	modpackData, err := ioutil.ReadFile(modpackFile)
	if err != nil {
		conn.Log("Unable to reload data " + err.Error())
		return
	}
	log.Println(string(modpackData))
	var modpacks []types.Modpack
	err = json.Unmarshal(modpackData, &modpacks)
	if err != nil {
		conn.Log("Could not parse json data " + err.Error() + "\n" + string(modpackData))
		return
	}
	for i, _ := range modpacks {
		modpack := &modpacks[i]
		modpack.Technic.Upload.AWS.SecretKey = encryption.DecryptString(modpack.Technic.Upload.AWS.SecretKey)
		modpack.Technic.Upload.AWS.AccessKey = encryption.DecryptString(modpack.Technic.Upload.AWS.AccessKey)
		modpack.Solder.Password = encryption.DecryptString(modpack.Solder.Password)
	}
	conn.Write("data-loaded", modpacks)
}
