package handlers

import (
	"io/ioutil"
	"log"
)

type inputDirData struct {
	InputDir string `json:"inputDir"`
}

func findAdditionalFolders(conn websocketConnection, data interface{}) {
	conn.Log("Test")
	dir := data.(inputDirData).InputDir
	files, _ := ioutil.ReadDir(dir)
	folders := []string{}
	// Iterate the files
	for _, file := range files {
		// We only need the directories
		if file.IsDir() {
			// The mods folder should be handles in a special way
			if file.Name() == "mods" {
				subFiles, _ := ioutil.ReadDir(file.Name())
				for _, subfile := range subFiles {
					if subfile.IsDir() {
						folders = append(folders, subfile.Name())
					}
				}
			} else {
				folders = append(folders, file.Name())
			}
		}
	}
	log.Println(folders)
}
