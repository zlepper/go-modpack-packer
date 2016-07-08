package types

import (
	"encoding/hex"
	"github.com/zlepper/go-modpack-packer/source/backend/helpers"
	"log"
	"regexp"
	"strings"
)

type Mod struct {
	ModId            string         `json:"modid"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	Version          string         `json:"version"`
	MinecraftVersion string         `json:"mcversion"`
	Url              string         `json:"url"`
	Authors          string         `json:"authors"`
	Credits          string         `json:"credits"`
	Filename         string         `json:"filename"`
	Md5              string         `json:"md5"`
	IsOnSolder       bool           `json:"isOnSolder"`
	Permission       *UserPermission `json:"userPermission,omitempty"`
}

func (m *Mod) GenerateOnlineVersion() string {
	return m.MinecraftVersion + "-" + m.Version
}

func (m *Mod) NormalizeId() {
	reg := []rune("\\\\|\\/|\\||:|\\*|\\\"|<|>|'|\\?|&|\\$|@|=|;|\\+|\\s|,|{|}|\\^|%|`|\\]|\\[|~|#|_") // Also known as the Fuck You Regex
	for i := 0; i < 32; i++ {
		c := rune(i)
		reg = append(reg, '|')
		reg = append(reg, c)
	}
	for i := 127; i < 256; i++ {
		c := rune(i)
		reg = append(reg, '|')
		reg = append(reg, c)
	}
	re := regexp.MustCompile(string(reg))
	m.ModId = strings.ToLower(re.ReplaceAllString(m.ModId, ""))
}

func (m *Mod) GetVersionString() string {
	return m.ModId + "-" + m.MinecraftVersion + "-" + m.Version
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
