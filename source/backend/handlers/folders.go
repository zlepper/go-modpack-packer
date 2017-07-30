package handlers

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

func getFolderList(topFolder string) ([]string, error) {
	infos, err := ioutil.ReadDir(topFolder)
	log.Println(err)
	if err != nil {
		// Incase the specified file doesn't exist, we'll just show the level above
		if os.IsNotExist(err) {
			log.Println("getting super folders")
			return getFolderList(path.Dir(topFolder))
		}
		return []string{}, err
	}

	var filenames []string
	for _, info := range infos {
		if info.IsDir() {
			filenames = append(filenames, path.Join(topFolder, info.Name()))
		}
	}

	return filenames, nil
}
