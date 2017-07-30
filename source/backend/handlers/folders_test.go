package handlers

import (
	"log"
	"testing"
)

func TestGetFolder(t *testing.T) {
	folders, err := getFolderList("C:/")
	if err != nil {
		t.Error(err)
	}

	log.Println(folders)
}
