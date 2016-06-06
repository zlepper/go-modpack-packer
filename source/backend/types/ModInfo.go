package types

import "strings"

type ModInfo struct {
	ModId            string    `json:"modid"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Version          string    `json:"version"`
	MinecraftVersion string    `json:"mcversion"`
	Url              string    `json:"url"`
	AuthorList       []string  `json:"authorList"`
	Authors          []string  `json:"authors"`
	Author           string    `json:"author"`
	Credits          string    `json:"credits"`
	ModListVersion   int32     `json:"modListVersion"`
	ModList          []ModInfo `json:"modList"`
}

func (m *ModInfo) CreateModResponse(filename string) Mod {
	modRes := Mod{
		ModId:            strings.ToLower(m.ModId),
		Name:             m.Name,
		Description:      m.Description,
		MinecraftVersion: m.MinecraftVersion,
		Url:              m.Url,
		Version:          m.Version,
		Credits:          m.Credits,
		Filename:         filename,
	}
	if len(m.AuthorList) > 0 {
		modRes.Authors = strings.Join(m.AuthorList, ", ")
	} else if len(m.Authors) > 0 {
		modRes.Authors = strings.Join(m.Authors, ", ")
	} else if m.Author != "" {
		modRes.Authors = m.Author
	}
	return modRes
}
