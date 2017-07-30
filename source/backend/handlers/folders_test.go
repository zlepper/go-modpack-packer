package handlers

import (
	"testing"
	"log"
)

func TestGetFolder(t *testing.T) {
	folders, err := getFolderList("C:/")
	if err != nil {
		t.Error(err)
	}

	log.Println(folders)
}
