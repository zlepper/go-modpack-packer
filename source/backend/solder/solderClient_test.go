package solder

import "testing"

func TestSolderClient_GetModId(t *testing.T) {
	client := NewSolderClient("http://solder.zlepper.dk/index.php")
	client.Login("testuser", "password")

	id := client.GetModId("agricraft")
	if id != "508" {
		t.Errorf("Expected mod id 508, got %s", id)
	}
}
