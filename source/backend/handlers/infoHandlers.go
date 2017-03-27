package handlers

import (
	"archive/zip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/mitchellh/mapstructure"
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
	"time"
)

var checkPermissions bool
var checkPublicPermissions bool

func gatherInformation(conn websocket.WebsocketConnection, data interface{}) {
	modpack := types.CreateSingleModpackData(data)
	checkPermissions = modpack.Technic.CheckPermissions
	checkPublicPermissions = modpack.Technic.IsPublicPack
	gatherInformationAboutMods(path.Join(modpack.InputDirectory, "mods"), conn)
}

func gatherInformationAboutMods(inputDirectory string, conn websocket.WebsocketConnection) {
	t1 := time.Now()
	filesAndDirectories, err := ioutil.ReadDir(inputDirectory)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	files := make([]os.FileInfo, 0)
	for _, f := range filesAndDirectories {
		if f.IsDir() {
			continue
		}
		// Files starting with . (dot) er hidden files under both OSX and linux
		// and they should be ignored
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		files = append(files, f)
	}
	fmt.Println(len(files))
	var waiter sync.WaitGroup
	conn.Write("total-mod-files", len(files))
	for _, f := range files {
		waiter.Add(1)
		fullname := path.Join(inputDirectory, f.Name())
		go gatherInformationAboutMod(fullname, conn, &waiter)
	}
	waiter.Wait()
	conn.Write("all-mod-files-scanned", "")
	t2 := time.Since(t1).Nanoseconds()
	fmt.Printf("Exploration time: %d ns\n", t2)
}

func gatherInformationAboutMod(modfile string, conn websocket.WebsocketConnection, waitGroup *sync.WaitGroup) {
	// Check if we already have the mod in the database. If we do we should just send that data to the client
	// instead of working through the zip file and calculating everything again.
	md5, err := helpers.ComputeMd5(modfile)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Println(err)
		return
	}
	md5String := hex.EncodeToString(md5)
	possibleMod := db.GetModsDb().GetModFromMd5(md5String)
	if possibleMod != nil {
		possibleMod.Filename = modfile
		sendModDataReady(*possibleMod, conn)
		waitGroup.Done()
		return
	}

	// The mod was not in the database, so time for some data crunching
	reader, err := zip.OpenReader(modfile)
	if err != nil {
		if err == zip.ErrFormat {
			conn.Log("err: " + modfile + " is not a valid zip file")
			sendModDataReady(types.Mod{Filename: modfile}, conn)
			waitGroup.Done()
			return
		} else {
			raven.CaptureErrorAndWait(err, nil)
			log.Panic(err)
		}
	}
	defer reader.Close()

	// Iterate the files in the archive to find the info files
	var foundInfoFile bool
	for _, f := range reader.File {
		// We only need .info and a certain .json file
		if (strings.HasSuffix(f.Name, ".info") && strings.Index(f.Name, "dependancies") == -1 && strings.Index(f.Name, "dependencies") == -1) || f.Name == "litemod.json" {
			r, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			readInfoFile(r, conn, f.FileInfo().Size(), modfile)
			foundInfoFile = true
			r.Close()
		}
	}
	if !foundInfoFile {
		sendModDataReady(types.Mod{Filename: modfile}, conn)
	}
	waitGroup.Done()
}

func readInfoFile(file io.ReadCloser, conn websocket.WebsocketConnection, size int64, filename string) {
	content := make([]byte, size)
	_, err := file.Read(content)
	content = []byte(strings.Replace(strings.Replace(string(content), "\n", " ", -1), "\r", "", -1))
	if err != nil {
		// For some reason the zip file reader in GO 1.7 gives io.EOF when reaching the end of the
		// file, which means the file.read will return an error, even though it read the content successfully...
		// Because why the f**k not?!
		if err != io.EOF {
			raven.CaptureError(err, nil)
			conn.Log(err.Error() + "\n" + filename + "\n" + string(debug.Stack()))
			return
		}
	}
	var mod types.ModInfo
	normalMod := make([]types.ModInfo, 0)
	err = json.Unmarshal(content, &normalMod)
	if err != nil {
		// Try with mod version 2, or with litemod
		err = json.Unmarshal(content, &mod)
		if err != nil {
			raven.CaptureError(err, nil)
			conn.Log(err.Error() + "\n" + string(content) + "\n" + filename)
			sendModDataReady(types.Mod{Filename: filename}, conn)
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
			modRes = *possibleMatch
		}
	}
	sendModDataReady(modRes, conn)
}

const modDataReadyEvent string = "mod-data-ready"

func sendModDataReady(mod types.Mod, conn websocket.WebsocketConnection) {
	if checkPermissions {
		permissionsDb := db.GetPermissionsDb()
		permission := permissionsDb.GetPermissionPolicy(mod.ModId, checkPublicPermissions)
		if permission != types.Open {
			mod.Permission = db.GetModsDb().GetModPermission(mod.ModId)
			if mod.Permission == nil {
				mod.Permission = &types.UserPermission{}
			}
			mod.Permission.Policy = permission
			if permission != types.Unknown && mod.Permission.PermissionLink == "" {
				// We have some data on this mod, might as well fill it in
				permissionData := permissionsDb.GetPermission(mod.ModId)
				mod.Permission.LicenseLink = permissionData.LicenseLink
				mod.Permission.ModLink = permissionData.ModLink
			}
		} else {
			entry := permissionsDb.GetPermission(mod.ModId)
			mod.Permission = &types.UserPermission{
				Policy:      permission,
				LicenseLink: entry.LicenseLink,
				ModLink:     entry.ModLink,
			}
		}
		log.Printf("Mod '%s' has permission '%v'", mod.ModId, *mod.Permission)
	}
	mod.NormalizeAll()
	conn.Write(modDataReadyEvent, mod)
}

const gotPermissionDataEvent string = "got-permission-data"

func CheckPermissionStore(conn websocket.WebsocketConnection, data interface{}) {
	type dataSearch struct {
		ModId    string `json:"modId"`
		IsPublic bool   `json:"isPublic"`
	}

	var search dataSearch
	mapstructure.Decode(data, &search)
	// First check own database
	modsDb := db.GetModsDb()
	permissions := modsDb.GetModPermission(search.ModId)
	if permissions != nil {
		permissions.ModId = search.ModId
		conn.Write(gotPermissionDataEvent, permissions)
		return
	}

	// We didn't find anything, so lets check the permissions store
	permissionsDb := db.GetPermissionsDb()
	p := permissionsDb.GetPermission(search.ModId)
	if p == nil {
		conn.Write(gotPermissionDataEvent, types.UserPermission{
			Policy: types.Unknown,
			ModId:  search.ModId,
		})
		return
	}

	var policy types.PermissionPolicy
	if search.IsPublic {
		policy = p.PublicPolicy
	} else {
		policy = p.PrivatePolicy
	}

	permissions = &types.UserPermission{
		Policy:      policy,
		LicenseLink: p.LicenseLink,
		ModLink:     p.ModLink,
		ModId:       search.ModId,
	}
	conn.Write(gotPermissionDataEvent, permissions)

}
