package types

import "testing"

type data struct {
	mod    Mod
	result bool
}

var tests = []data{
	{Mod{
		ModId:            "test",
		Name:             "Test mod",
		Version:          "1.1.2",
		MinecraftVersion: "1.7.10",
		Authors:          "TestGuy",
	}, true},
	{Mod{
		ModId:            "",
		Name:             "Test mod",
		Version:          "1.1.2",
		MinecraftVersion: "1.7.10",
		Authors:          "TestGuy",
	}, false},
	{Mod{
		ModId:            "test",
		Name:             "",
		Version:          "1.1.2",
		MinecraftVersion: "1.7.10",
		Authors:          "TestGuy",
	}, false},
	{Mod{
		ModId:            "test",
		Name:             "Test mod",
		Version:          "",
		MinecraftVersion: "1.7.10",
		Authors:          "TestGuy",
	}, false},
	{Mod{
		ModId:            "test",
		Name:             "Test mod",
		Version:          "1.1.2",
		MinecraftVersion: "",
		Authors:          "TestGuy",
	}, false},
	{Mod{
		ModId:            "test",
		Name:             "Test mod",
		Version:          "1.1.2",
		MinecraftVersion: "1.7.10",
		Authors:          "",
	}, false},
	{Mod{
		ModId:            "example",
		Name:             "Test mod",
		Version:          "1.1.2",
		MinecraftVersion: "1.7.10",
		Authors:          "TestGuy",
	}, false},
	{Mod{
		ModId:            "test",
		Name:             "example",
		Version:          "1.1.2",
		MinecraftVersion: "1.7.10",
		Authors:          "TestGuy",
	}, false},
	{Mod{
		ModId:            "test",
		Name:             "Test mod",
		Version:          "example",
		MinecraftVersion: "1.7.10",
		Authors:          "TestGuy",
	}, false},
}

func TestIsValid(t *testing.T) {
	for _, d := range tests {
		valid := d.mod.IsValid()
		if valid != d.result {
			t.Error("Expected", d.result, "got", valid)
		}
	}
}
