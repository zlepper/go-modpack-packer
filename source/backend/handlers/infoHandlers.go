package handlers

import (
	"io/ioutil"
	"os"
	"path"
	"archive/zip"
	"log"
	"strings"
	"io"
	"encoding/json"
	"runtime/debug"
	"regexp"
)

func gatherInformation(conn websocketConnection, data interface{}) {
	modpack := createSingleModpackData(data)
	gatherInformationAboutMods(path.Join(modpack.InputDirectory, "mods"), conn)
}

func gatherInformationAboutMods(inputDirectory string, conn websocketConnection) {
	filesAndDirectories, _ := ioutil.ReadDir(inputDirectory)
	files := make([]os.FileInfo, 0)
	for _, f := range filesAndDirectories {
		if f.IsDir() {
			continue
		}
		files = append(files, f)
	}
	Write(conn, "total-mod-files", len(files))
	for _, f := range files {
		fullname := path.Join(inputDirectory, f.Name())
		go gatherInformationAboutMod(fullname, conn)
	}
}

func gatherInformationAboutMod(modfile string, conn websocketConnection) {
	reader, err := zip.OpenReader(modfile)
	if err != nil {
		if err == zip.ErrFormat {
			conn.Log("err: " + modfile + " is not a valid zip file")
			return
		} else {
			log.Panic(err)
		}
	}
	defer reader.Close()

	// Iterate the files in the archive to find the info files
	for _, f := range reader.File {
		// We only need .info and a certain .json file
		if strings.HasSuffix(f.Name, "mod.info") || f.Name == "litemod.json" {
			r, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			readInfoFile(r, conn, f.FileInfo().Size(), modfile)
		}
	}
}

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

func (m *Mod) normalizeId() {
	reg := []rune("\\\\|\\/|\\||:|\\*|\\\"|<|>|'|\\?|&|\\$|@|=|;|\\+|\\s|,|{|}|\\^|%|`|\\]|\\[|~|#") // Also known as the Fuck You Regex
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
	m.ModId = re.ReplaceAllString(m.ModId, "")
}

type ModInfo struct {
	ModId            string `json:"modid"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Version          string `json:"version"`
	MinecraftVersion string `json:"mcversion"`
	Url              string `json:"url"`
	AuthorList       []string `json:"authorList"`
	Authors          []string `json:"authors"`
	Author           string `json:"author"`
	Credits          string `json:"credits"`
	ModListVersion   int32 `json:"modListVersion"`
	ModList          []ModInfo `json:"modList"`
}

func (m *ModInfo) createModResponse(filename string) Mod {
	modRes := Mod{
		ModId:m.ModId,
		Name:m.Name,
		Description:m.Description,
		MinecraftVersion:m.MinecraftVersion,
		Url:m.Url,
		Version:m.Version,
		Credits:m.Credits,
		Filename:filename,
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

func (m *Mod) getVersionString() string {
	return m.ModId + "-" + m.MinecraftVersion + "-" + m.Version
}

func readInfoFile(file io.ReadCloser, conn websocketConnection, size int64, filename string) {
	content := make([]byte, size)
	_, err := file.Read(content)
	content = []byte(strings.Replace(string(content), "\n", " ", -1))
	if err != nil {
		conn.Log(err.Error() + "\n" + string(debug.Stack()))
	}
	var mod ModInfo
	normalMod := make([]ModInfo, 0)
	err = json.Unmarshal(content, &normalMod)
	if err != nil {
		// Try with mod version 2, or with litemod
		err = json.Unmarshal(content, &mod)
		if err != nil {
			conn.Log(err.Error() + "\n" + string(content) + "\n" + filename)
		}
		// Handle version 2 mods
		if mod.ModListVersion == 2 {
			createModResponse(conn, mod.ModList[0], filename)
		} else {
			// Handle liteloader mods
			createModResponse(conn, mod, filename)
		}
		return
	}
	if len(normalMod) > 0 {
		createModResponse(conn, normalMod[0], filename)
	} else {
		createModResponse(conn, ModInfo{}, filename)
	}

}

func createModResponse(conn websocketConnection, mod ModInfo, filename string) {
	const modDataReadyEvent string = "mod-data-ready"
	modRes := mod.createModResponse(filename)
	modRes.normalizeId()
	Write(conn, modDataReadyEvent, modRes)
}
