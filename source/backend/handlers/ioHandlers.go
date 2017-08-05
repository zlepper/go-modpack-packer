package handlers

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"github.com/mitchellh/mapstructure"
	"github.com/zlepper/go-modpack-packer/source/backend/consts"
	"github.com/zlepper/go-modpack-packer/source/backend/encryption"
	"github.com/zlepper/go-modpack-packer/source/backend/internal"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

var mutex *sync.Mutex

func init() {
	mutex = &sync.Mutex{}
}

type inputDirData struct {
	InputDir string `json:"inputDir"`
	Key      int    `json:"key"`
}

type foundFolders struct {
	Folders []string `json:"folders"`
	Key     int      `json:"key"`
}

func findAdditionalFolders(conn types.WebsocketConnection, data interface{}) {

	var dir inputDirData
	err := mapstructure.Decode(data, &dir)
	if err != nil {
		log.Panicln(err)
	}

	inputDir := dir.InputDir

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	folders := []string{}
	// Iterate the files
	for _, file := range files {
		// We only need the directories
		if file.IsDir() {
			// The mods folder should be handled in a special way
			if file.Name() == "mods" {
				subFiles, _ := ioutil.ReadDir(filepath.Join(inputDir, file.Name()))
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
	conn.Write("found-folders", foundFolders{Folders: folders, Key: dir.Key})
}

func saveModpacks(conn types.WebsocketConnection, data interface{}) {
	internal.OutstandingProcess.Add(1)
	// Get the appData directory, since go doesn't expose it, electron passes it as a parameter
	dataDirectory := consts.DataDirectory
	modpackFile := filepath.Join(dataDirectory, "modpacks.json")
	modpacks := types.CreateModpackData(data)
	for i := range modpacks {
		modpack := &modpacks[i]
		modpack.Technic.Upload.AWS.SecretKey = encryption.EncryptString(modpack.Technic.Upload.AWS.SecretKey)
		modpack.Technic.Upload.AWS.AccessKey = encryption.EncryptString(modpack.Technic.Upload.AWS.AccessKey)
		modpack.Technic.Upload.FTP.Password = encryption.EncryptString(modpack.Technic.Upload.FTP.Password)
		modpack.Technic.Upload.FTP.Username = encryption.EncryptString(modpack.Technic.Upload.FTP.Username)
		modpack.Solder.Password = encryption.EncryptString(modpack.Solder.Password)
	}
	mutex.Lock()
	f, err := os.Create(modpackFile)
	defer f.Close()
	if err != nil {
		mutex.Unlock()
		log.Println(err)
		conn.Error(err.Error())
		raven.CaptureError(err, nil)
		return
	}
	err = json.NewEncoder(f).Encode(modpacks)
	if err != nil {
		mutex.Unlock()
		log.Println(err)
		conn.Error(err.Error())
		raven.CaptureError(err, nil)
		return
	}
	mutex.Unlock()
	internal.OutstandingProcess.Done()
}

func loadModpacks(conn types.WebsocketConnection) {
	dataDirectory := consts.DataDirectory
	modpackFile := filepath.Join(dataDirectory, "modpacks.json")
	mutex.Lock()
	modpackData, err := ioutil.ReadFile(modpackFile)
	mutex.Unlock()
	var modpacks []types.Modpack = make([]types.Modpack, 0)
	if err != nil {
		log.Println("Unable to reload data")
		log.Println(err)
		conn.Log("Unable to reload data " + err.Error())
		conn.Write("data-loaded", modpacks)
		return
	}
	err = json.Unmarshal(modpackData, &modpacks)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Println("Could not parse json data " + err.Error() + "\n" + string(modpackData))
		conn.Error("Could not parse json data. Please check the logs.")
		conn.Write("data-loaded", modpacks)
		return
	}
	for i := range modpacks {
		modpack := &modpacks[i]
		modpack.Technic.Upload.AWS.SecretKey = encryption.DecryptString(modpack.Technic.Upload.AWS.SecretKey)
		modpack.Technic.Upload.AWS.AccessKey = encryption.DecryptString(modpack.Technic.Upload.AWS.AccessKey)
		modpack.Technic.Upload.FTP.Password = encryption.DecryptString(modpack.Technic.Upload.FTP.Password)
		modpack.Technic.Upload.FTP.Username = encryption.DecryptString(modpack.Technic.Upload.FTP.Username)
		modpack.Solder.Password = encryption.DecryptString(modpack.Solder.Password)
	}
	conn.Write("data-loaded", modpacks)
}

type FolderRequest struct {
	Folder string `json:"folder"`
	Key    int    `json:"key"`
}

type FolderResponse struct {
	Folders []string `json:"folders"`
	Key     int      `json:"key"`
}

func getFolders(conn types.WebsocketConnection, fr interface{}) {

	var folders []string
	var err error

	var folderRequest FolderRequest
	err = mapstructure.Decode(fr, &folderRequest)
	if err != nil {
		conn.Error(err.Error())
		return
	}

	path := folderRequest.Folder

	log.Println("Getting folders for " + path)

	if path == "/" {
		folders, err = GetDrives()
	} else {
		folders, err = getFolderList(path)
	}

	if err != nil {
		conn.Error(err.Error())
	} else {
		if len(folders) == 0 {
			folders = append(folders, path)
		}
		sort.Strings(folders)
		log.Println(folders)
		conn.Write("got-folders", FolderResponse{Folders: folders, Key: folderRequest.Key})
	}
}
