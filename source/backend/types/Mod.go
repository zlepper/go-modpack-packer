package types

import (
	"regexp"
	"strings"
)

type Mod struct {
	ModId            string `json:"modid"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Version          string `json:"version"`
	MinecraftVersion string `json:"mcversion"`
	Url              string `json:"url"`
	Authors          string `json:"authors"`
	Credits          string `json:"credits"`
	Filename         string `json:"filename"`
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
