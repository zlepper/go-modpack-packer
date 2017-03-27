package types

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"regexp"
	"strings"
)

type ForgeVersion struct {
	Build            float64 `json:"build"`
	DownloadUrl      string  `json:"downloadUrl"`
	MinecraftVersion string  `json:"minecraftVersion"`
}

type AWSConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
}

type FtpConfig struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Path     string `json:"path"`
}

type UploadConfig struct {
	Type string    `json:"type"`
	AWS  AWSConfig `json:"aws"`
	FTP  FtpConfig `json:"ftp"`
}

type TechnicConfig struct {
	IsSolderPack     bool         `json:"isSolderPack"`
	CreateForgeZip   bool         `json:"createForgeZip"`
	ForgeVersion     ForgeVersion `json:"forgeVersion"`
	CheckPermissions bool         `json:"checkPermissions"`
	IsPublicPack     bool         `json:"isPublicPack"`
	Memory           float64      `json:"memory"`
	Java             string       `json:"java"`
	Upload           UploadConfig `json:"upload",mapstructure:"upload"`
	RepackAllMods    bool         `json:"repackAllMods"`
}

type FtbConfig struct {
	IsPublicPack bool `json:"isPublicPack"`
}

type Folder struct {
	Name    string `json:"name"`
	Include bool   `json:"include"`
}

type SolderInfo struct {
	Use      bool   `json:"use"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Modpack struct {
	Name                 string        `json:"name"`
	InputDirectory       string        `json:"inputDirectory"`
	OutputDirectory      string        `json:"outputDirectory"`
	ClearOutputDirectory bool          `json:"clearOutputDirectory"`
	MinecraftVersion     string        `json:"minecraftVersion"`
	Version              string        `json:"version"`
	AdditionalFolders    []Folder      `json:"additionalFolders"`
	Technic              TechnicConfig `json:"technic"`
	Ftb                  FtbConfig     `json:"ftb"`
	Solder               SolderInfo    `json:"solder"`
	IsNew                bool          `json:"isNew"`
}

func (m *Modpack) GetSlug() string {
	re := regexp.MustCompile("\\|/|\\||:|\\*|\"|<|>|\\?|'")
	s := re.ReplaceAllString(m.Name, "")
	s = strings.Replace(s, " ", "-", -1)
	s = strings.ToLower(s)
	return s
}

func (m *Modpack) GetVersionString() string {
	return m.MinecraftVersion + "-" + m.Version
}

func CreateSingleModpackData(di interface{}) Modpack {
	d := di.(map[string]interface{})
	//d := data.(map[string]interface{})
	var modpack Modpack
	err := mapstructure.Decode(d, &modpack)
	if err != nil {
		log.Panic(err)
	}
	return modpack
}

func CreateModpackData(data interface{}) []Modpack {
	dL := data.([]interface{})
	modpacks := make([]Modpack, 0)
	for _, di := range dL {
		modpack := CreateSingleModpackData(di)
		modpacks = append(modpacks, modpack)
	}
	return modpacks
}
