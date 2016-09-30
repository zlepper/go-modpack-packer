package solder

import (
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"testing"
)

func getClient() *SolderClient {
	client := NewSolderClient("http://solder.zlepper.dk/index.php")
	client.Login("test@test.com", "password")

	return client
}

func TestSolderClient_GetModId(t *testing.T) {
	client := getClient()

	id := client.GetModId("ae2stuff")
	if id != "771" {
		t.Errorf("Expected mod id 771, got %s", id)
	}

	id = client.GetModId("teknflganfgcraft")
	if id != "" {
		t.Error("Expected mod id to be empty, got", id)
	}
}

func TestSolderClient_GetModVersionId(t *testing.T) {
	client := getClient()

	modVersionId := client.GetModVersionId(&types.OutputInfo{
		Id:               "ae2stuff",
		MinecraftVersion: "1.7.10",
		Version:          "0.5.1.9",
	})

	if modVersionId != "1912" {
		t.Errorf("Expected mod version id 1912 got %s", modVersionId)
	}

	modVersionId = client.GetModVersionId(&types.OutputInfo{
		Id:               "ae2stuff",
		MinecraftVersion: "1.7.10",
		Version:          "1.4.0",
	})

	if modVersionId != "" {
		t.Error("Expected not to find a mod version id, got", modVersionId)
	}
}

func TestSolderClient_IsModversionOnline(t *testing.T) {
	client := getClient()

	exists := client.IsModversionOnline(&types.OutputInfo{
		Id:               "agricraft",
		MinecraftVersion: "1.7.10",
		Version:          "1.5.0",
	})

	if !exists {
		t.Error("Expected modversion to exists. It didn't")
	}

	exists = client.IsModversionOnline(&types.OutputInfo{
		Id:               "acricraft",
		MinecraftVersion: "1.7.10",
		Version:          "1.4.0",
	})

	if exists {
		t.Error("Expected modversion to not exists. It did")
	}
}
