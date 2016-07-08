package handlers

import (
	"archive/zip"
	"encoding/hex"
	"encoding/json"
	"github.com/zlepper/go-modpack-packer/source/backend/db"
	"github.com/zlepper/go-modpack-packer/source/backend/helpers"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"github.com/zlepper/go-websocket-connection"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"sync"
)

func gatherInformation(conn websocket.WebsocketConnection, data interface{}) {
	modpack := types.CreateSingleModpackData(data)
	gatherInformationAboutMods(path.Join(modpack.InputDirectory, "mods"), conn)
}

func gatherInformationAboutMods(inputDirectory string, conn websocket.WebsocketConnection) {
	filesAndDirectories, _ := ioutil.ReadDir(inputDirectory)
	files := make([]os.FileInfo, 0)
	for _, f := range filesAndDirectories {
		if f.IsDir() {
			continue
		}
		files = append(files, f)
	}
	var waiter sync.WaitGroup
	conn.Write("total-mod-files", len(files))
	for _, f := range files {
		waiter.Add(1)
		fullname := path.Join(inputDirectory, f.Name())
		go gatherInformationAboutMod(fullname, conn, &waiter)
	}
	waiter.Wait()
	conn.Write("all-mod-files-scanned", "")
}

func gatherInformationAboutMod(modfile string, conn websocket.WebsocketConnection, waitGroup *sync.WaitGroup) {
	// Check if we already have the mod in the database. If we do we should just send that data to the client
	// instead of working through the zip file and calculating everything again.
	md5, err := helpers.ComputeMd5(modfile)
	md5String := hex.EncodeToString(md5)
	possibleMod := db.GetModsDb().GetModFromMd5(md5String)
	if possibleMod != nil {
		sendModDataReady(*possibleMod, conn)
		waitGroup.Done()
		return
	}

	// The mod was not in the database, so time for some data crunching
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
	waitGroup.Done()
}

func readInfoFile(file io.ReadCloser, conn websocket.WebsocketConnection, size int64, filename string) {
	content := make([]byte, size)
	_, err := file.Read(content)
	content = []byte(strings.Replace(string(content), "\n", " ", -1))
	if err != nil {
		conn.Log(err.Error() + "\n" + string(debug.Stack()))
		return
	}
	var mod types.ModInfo
	normalMod := make([]types.ModInfo, 0)
	err = json.Unmarshal(content, &normalMod)
	if err != nil {
		// Try with mod version 2, or with litemod
		err = json.Unmarshal(content, &mod)
		if err != nil {
			conn.Log(err.Error() + "\n" + string(content) + "\n" + filename)
			return
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
		createModResponse(conn, types.ModInfo{}, filename)
	}

}

func createModResponse(conn websocket.WebsocketConnection, mod types.ModInfo, filename string) {
	modRes := mod.CreateModResponse(filename)
	modRes.NormalizeId()

	md5 := modRes.GetMd5()
	if !modRes.IsValid() {
		modsDb := db.GetModsDb()
		possibleMatch := modsDb.GetModFromMd5(md5)
		if possibleMatch != nil {
			//if possibleMatch.IsValid() {
			// Should be valid, but better safe than sorry
			modRes = *possibleMatch
			//} else {
			//	fmt.Println("database mod was not valid!")
			//}
		}
	}
	sendModDataReady(modRes, conn)
}

func sendModDataReady(mod types.Mod, conn websocket.WebsocketConnection) {
	const modDataReadyEvent string = "mod-data-ready"
	conn.Write(modDataReadyEvent, mod)
}
