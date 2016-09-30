package types

import (
	"encoding/hex"
	"github.com/zlepper/go-modpack-packer/source/backend/helpers"
	"log"
	"regexp"
	"strings"
)

type Mod struct {
	ModId            string          `json:"modid"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	Version          string          `json:"version"`
	MinecraftVersion string          `json:"mcversion"`
	Url              string          `json:"url"`
	Authors          string          `json:"authors"`
	Credits          string          `json:"credits"`
	Filename         string          `json:"filename"`
	Md5              string          `json:"md5"`
	IsOnSolder       bool            `json:"isOnSolder"`
	Permission       *UserPermission `json:"userPermission,omitempty"`
}

func (m *Mod) GenerateOnlineVersion() string {
	return m.MinecraftVersion + "-" + m.Version
}

const (
	validChars = `[^\w\d-_ '"\./:` + "`]"
)

var re = regexp.MustCompile(validChars)

func normalize(s string) string {
	return re.ReplaceAllString(s, "")
}

func (m *Mod) NormalizeId() {
	m.ModId = strings.ToLower(strings.Replace(strings.Replace(normalize(m.ModId), ".", "", -1), " ", "-", -1))
}

func (m *Mod) GetVersionString() string {
	return strings.Replace(m.ModId+"-"+m.MinecraftVersion+"-"+m.Version, " ", "-", -1)
}

func (m *Mod) GetMd5() string {
	if m.Md5 != "" {
		return m.Md5
	}

	md5, err := helpers.ComputeMd5(m.Filename)
	if err != nil {
		log.Panic(err)
	}
	m.Md5 = hex.EncodeToString(md5)
	return m.Md5
}

func (m *Mod) IsValid() bool {
	return m.ModId != "" &&
		m.Name != "" &&
		m.Version != "" &&
		m.MinecraftVersion != "" &&
		m.Authors != "" &&
		!helpers.IgnoreCaseContains(m.ModId, "example") &&
		!helpers.IgnoreCaseContains(m.Name, "example") &&
		!helpers.IgnoreCaseContains(m.Version, "example") &&
		!helpers.IgnoreCaseContains(m.Name, "${") &&
		!helpers.IgnoreCaseContains(m.Version, "${") &&
		!helpers.IgnoreCaseContains(m.MinecraftVersion, "${") &&
		!helpers.IgnoreCaseContains(m.ModId, "${") &&
		!helpers.IgnoreCaseContains(m.Version, "@version@")
}

func (m *Mod) NormalizeAll() {
	m.NormalizeId()
	m.Version = strings.Replace(normalize(m.Version), " ", "-", -1)
	m.MinecraftVersion = normalize(m.MinecraftVersion)
	m.Name = normalize(m.Name)
	m.Description = normalize(m.Description)
	m.Url = normalize(m.Url)
	m.Authors = normalize(m.Authors)

	// Special check to remove minecraft version in mod version strings
	if m.Version != m.MinecraftVersion {
		if strings.Contains(m.Version, m.MinecraftVersion) {
			m.Version = strings.Replace(m.Version, m.MinecraftVersion, "", -1)
			m.Version = strings.Trim(m.Version, " -_.+")
		}
	}
}

func (m *Mod) SetSolderStatus(status bool) {
	m.IsOnSolder = status
}

func SafeNormalizeString(s string) string {
	s = strings.Replace(strings.ToLower(s), " ", "-", -1)
	return strings.Replace(s, ".", "", -1)
}

func (m *Mod) GenerateSimpleOutputInfo() *OutputInfo {
	return &OutputInfo{
		Name:             m.Name,
		Id:               SafeNormalizeString(m.ModId),
		Version:          m.Version,
		MinecraftVersion: m.MinecraftVersion,
		Description:      m.Description,
		Author:           m.Authors,
		ProgressKey:      m.Filename,
		IsOnSolder:       m.IsOnSolder,
		Permissions:      m.Permission,
	}
}
