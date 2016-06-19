package handlers

import (
	"encoding/json"
	"github.com/zlepper/go-modpack-packer/source/backend/encryption"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"github.com/zlepper/go-websocket-connection"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var mutex *sync.Mutex

func init() {
	mutex = &sync.Mutex{}
}

type inputDirData struct {
	InputDir string `json:"inputDir"`
}

func createInputDirData(data map[string]interface{}) inputDirData {
	var res inputDirData
	res.InputDir = data["inputDir"].(string)
	return res
}

func findAdditionalFolders(conn websocket.WebsocketConnection, data interface{}) {
	dir := createInputDirData(data.(map[string]interface{}))
	files, _ := ioutil.ReadDir(dir.InputDir)
	folders := []string{}
	// Iterate the files
	for _, file := range files {
		// We only need the directories
		if file.IsDir() {
			// The mods folder should be handled in a special way
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

func saveModpacks(conn websocket.WebsocketConnection, data interface{}) {
	// Get the appData directory, since go doesn't expose it, electron passes it as a parameter
	dataDirectory := os.Args[1]
	modpackFile := filepath.Join(dataDirectory, "modpacks.json")
	modpacks := types.CreateModpackData(data)
	for i, _ := range modpacks {
		modpack := &modpacks[i]
		modpack.Technic.Upload.AWS.SecretKey = encryption.EncryptString(modpack.Technic.Upload.AWS.SecretKey)
		modpack.Technic.Upload.AWS.AccessKey = encryption.EncryptString(modpack.Technic.Upload.AWS.AccessKey)
		modpack.Solder.Password = encryption.EncryptString(modpack.Solder.Password)
	}
	mutex.Lock()
	f, err := os.Create(modpackFile)
	defer f.Close()
	if err != nil {
		mutex.Unlock()
		log.Panic(err)
	}
	err = json.NewEncoder(f).Encode(modpacks)
	if err != nil {
		mutex.Unlock()
		log.Panic(err)
	}
	mutex.Unlock()
}

func loadModpacks(conn websocket.WebsocketConnection) {
	dataDirectory := os.Args[1]
	modpackFile := filepath.Join(dataDirectory, "modpacks.json")
	mutex.Lock()
	modpackData, err := ioutil.ReadFile(modpackFile)
	mutex.Unlock()
	var modpacks []types.Modpack = make([]types.Modpack, 0)
	if err != nil {
		conn.Log("Unable to reload data " + err.Error())
		conn.Write("data-loaded", modpacks)
		return
	}
	log.Println(string(modpackData))
	err = json.Unmarshal(modpackData, &modpacks)
	if err != nil {
		panic("Could not parse json data " + err.Error() + "\n" + string(modpackData))
	}
	for i, _ := range modpacks {
		modpack := &modpacks[i]
		modpack.Technic.Upload.AWS.SecretKey = encryption.DecryptString(modpack.Technic.Upload.AWS.SecretKey)
		modpack.Technic.Upload.AWS.AccessKey = encryption.DecryptString(modpack.Technic.Upload.AWS.AccessKey)
		modpack.Solder.Password = encryption.DecryptString(modpack.Solder.Password)
	}
	conn.Write("data-loaded", modpacks)
}
