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
	"github.com/zlepper/go-modpack-packer/source/backend/internal"
	"github.com/zlepper/go-modpack-packer/source/backend/solder"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/upload"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const donePackingPartName string = "done-packing-part"
const packingPartName string = "packing-part"

func build(conn types.WebsocketConnection, data interface{}) {
	dat := data.(map[string]interface{})
	modpack := types.CreateSingleModpackData(dat["modpack"])
	mods := make([]*types.Mod, 0)
	modsData := dat["mods"]
	modsDat, _ := json.Marshal(modsData)
	err := json.Unmarshal(modsDat, &mods)
	if err != nil {
		conn.Log("Could not create mod list")
		return
	}
	buildModpack(modpack, mods, conn)
}

type uploadWaiting struct {
	Modpack types.Modpack       `json:"modpack"`
	Infos   []*types.OutputInfo `json:"infos"`
}

func continueRunning(conn types.WebsocketConnection, data interface{}) {
	dict := data.(map[string]interface{})

	var uploadInfo uploadWaiting
	err := mapstructure.Decode(dict, &uploadInfo)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Panic(err)
	}

	solderclient, buildId := updateSolder(uploadInfo.Modpack, conn)
	conn.Write(solderCurrentlyDoingEvent, "Updating mods")
	var wg sync.WaitGroup
	for _, info := range uploadInfo.Infos {
		wg.Add(1)
		go func() {
			addInfoToSolder(info, buildId, conn, solderclient)
			wg.Done()
		}()
	}
	wg.Wait()
	conn.Write("done-updating", "")
}

func buildModpack(modpack types.Modpack, mods []*types.Mod, conn types.WebsocketConnection) {
	// Create output directory
	outputDirectory := path.Join(modpack.OutputDirectory, modpack.Name)
	os.MkdirAll(outputDirectory, os.ModePerm)
	var total int

	var ch chan *types.OutputInfo
	ch = make(chan *types.OutputInfo)

	startTime := time.Now()
	var solderClient *solder.SolderClient
	// Only create the solder client if we actually have to use solder
	if modpack.Solder.Use {
		solderClient = solder.NewSolderClient(modpack.Solder.Url)
		solderClient.Login(modpack.Solder.Username, modpack.Solder.Password)
	}

	// Handle forge
	if modpack.Technic.CreateForgeZip {
		total++
		go packForgeFolder(modpack, conn, outputDirectory, &ch, solderClient)
	}

	// Handle any additional folder
	for _, folder := range modpack.AdditionalFolders {
		if folder.Include {
			total++
			go packAdditionalFolder(modpack, folder.Name, outputDirectory, conn, &ch)
		}
	}

	infos := make([]*types.OutputInfo, 0)
	var wg sync.WaitGroup
	var lock sync.Mutex
	// Handle mods
	wg.Add(len(mods))
	for _, mod := range mods {
		go func(m *types.Mod) {
			mod.NormalizeAll()
			// If the mod already is on solder, then we should likely skip it
			// however the user can override this. If they do we should still pack all files
			if !modpack.Technic.RepackAllMods && solder.IsOnSolder(solderClient, m) {
				conn.Write(packingPartName, m.Filename)
				conn.Write(donePackingPartName, m.Filename)
				lock.Lock()
				infos = append(infos, GenerateOutputInfo(m, ""))
				lock.Unlock()
			} else {
				go packMod(m, conn, outputDirectory, &ch)
			}
			wg.Done()
		}(mod)
		total++
	}
	conn.Write("total-to-pack", total)
	wg.Wait()

	// Save the mods to the database
	d := db.GetModsDb()
	for _, mod := range mods {
		mod.SetSolderStatus(true)
		d.AddMod(mod)
	}
	d.Save()

	count := len(infos)
	for count < total {
		info := <-ch
		infos = append(infos, info)
		count++
	}
	endTime := time.Now()
	spendtime := endTime.UnixNano() - startTime.UnixNano()
	spendtime = spendtime / (int64(time.Millisecond) / int64(time.Nanosecond))
	fmt.Printf("Time spend packing: %d ms", spendtime)

	switch modpack.Technic.Upload.Type {
	case "none":
		{
			// Send back information and await clearance from client
			conn.Write("waiting-for-file-upload", uploadWaiting{Modpack: modpack, Infos: infos})
			return
		}
	case "ftp":
		{
			upload.UploadFilesToFtp(&modpack, infos, conn)
		}
	case "s3":
		{
			upload.UploadFilesToS3(&modpack, infos, conn)
		}
	case "gfs":
		{
			upload.UploadFilesToGfs(&modpack, infos, conn)
		}
	}

	if modpack.Solder.Use {
		solderclient, buildId := updateSolder(modpack, conn)

		conn.Write(solderCurrentlyDoingEvent, "Updating mods")
		var wg sync.WaitGroup
		for _, info := range infos {
			localInfo := *info
			log.Println(localInfo)
			wg.Add(1)
			go func() {
				addInfoToSolder(&localInfo, buildId, conn, solderclient)
				wg.Done()
			}()
		}
		wg.Wait()
		conn.Write("done-updating", "")
	}

}

const solderCurrentlyDoingEvent string = "solder-currently-doing"

func updateSolder(modpack types.Modpack, conn types.WebsocketConnection) (*solder.SolderClient, string) {
	// Create solder client
	conn.Write(solderCurrentlyDoingEvent, "Logging in to solder")
	solderclient := solder.NewSolderClient(modpack.Solder.Url)
	err := solderclient.Login(modpack.Solder.Username, modpack.Solder.Password)
	if err != nil {
		log.Panic("Could not login to solder with the supplied credentials.", err)
	}

	conn.Write(solderCurrentlyDoingEvent, "Checking if modpack already exists")
	var modpackId string
	if solderclient.IsPackOnline(&modpack) {
		conn.Write(solderCurrentlyDoingEvent, "Modpack exists, getting info")
		modpackId = solderclient.GetModpackId(modpack.GetSlug())
	} else {
		conn.Write(solderCurrentlyDoingEvent, "Modpack does not exist. Creating...")
		modpackId = solderclient.CreatePack(modpack.Name, modpack.GetSlug())
	}

	var buildId string
	conn.Write(solderCurrentlyDoingEvent, "Checking if build exists")
	if solderclient.IsBuildOnline(&modpack) {
		conn.Write(solderCurrentlyDoingEvent, "Build exists, getting info")
		buildId = solderclient.GetBuildId(&modpack)
	} else {
		conn.Write(solderCurrentlyDoingEvent, "Build does not exist. Creating..")
		buildId = solderclient.CreateBuild(&modpack, modpackId)
	}
	return solderclient, buildId
}

func addInfoToSolder(info *types.OutputInfo, buildId string, conn types.WebsocketConnection, solderclient *solder.SolderClient) {
	conn.Write("updating-solder", info.ProgressKey)
	var modid string
	modid = solderclient.GetModId(info.Id)
	if modid == "" {
		log.Println("Could not get mod id. Adding mod to solder " + info.Id)
		modid = solderclient.AddMod(info)
	}
	if modid == "" {
		log.Println("Something went wrong when adding a mod to solder.")
		log.Printf("%v\n", *info)
		log.Printf("Application version: %s\n", internal.Version)
		log.Panic("Error. See above lines")
	}

	if info.File != "" {
		md5, err := helpers.ComputeMd5(info.File)
		if !solderclient.IsModversionOnline(info) {
			if err != nil {
				raven.CaptureError(err, nil)
				log.Panic(err)
			}
			log.Println("Adding mod version to solder for " + info.Name)
			solderclient.AddModVersion(modid, hex.EncodeToString(md5), info.GenerateOnlineVersion())
		} else {
			log.Println("Rehashing mod version for " + info.Name)
			id := solderclient.GetModVersionId(info)
			solderclient.RehashModVersion(id, hex.EncodeToString(md5))
		}
		go db.GetModsDb().MarkModAsOnSolder(hex.EncodeToString(md5))
	}
	if !solderclient.IsModversionActiveInBuild(info, buildId) {
		if solderclient.IsModInBuild(info, buildId) {
			log.Println("Mod is already in build, updating")
			solderclient.SetModVersionInBuild(info, buildId)
		} else {
			log.Println("Mod is not in build")
			solderclient.AddModversionToBuild(info, buildId)
		}
	}
	conn.Write("done-updating-solder", info.ProgressKey)
}

func packForgeFolder(modpack types.Modpack, conn types.WebsocketConnection, outputDirectory string, ch *chan *types.OutputInfo, sc *solder.SolderClient) {
	const minecraftForge string = "Minecraft Forge"
	version := fmt.Sprintf("%v", modpack.Technic.ForgeVersion.Build)
	conn.Write(packingPartName, minecraftForge)
	isOnSolder := false
	if sc != nil {
		isOnSolder = solder.IsOnSolder(sc, &types.Mod{
			Name:             minecraftForge,
			ModId:            "forge",
			Version:          version,
			MinecraftVersion: modpack.MinecraftVersion,
			Description:      "The core of everything modded minecraft.",
			Url:              "http://www.minecraftforge.net/",
			Authors:          "LexManos, cpw",
		})
	}
	var outputFile string
	if !isOnSolder {
		outputDirectory = path.Join(outputDirectory, "mods", "forge")
		os.MkdirAll(outputDirectory, os.ModePerm)
		outputFile = path.Join(outputDirectory, "forge-"+modpack.MinecraftVersion+"-"+version+".zip")

		zipfile, err := os.Create(outputFile)
		if err != nil {
			raven.CaptureError(err, nil)
			conn.Log("Error when creating zip file: " + err.Error() + "\n" + outputFile)
			return
		}
		defer zipfile.Close()

		zipWriter := zip.NewWriter(zipfile)
		defer zipWriter.Close()

		resp, err := http.Get(modpack.Technic.ForgeVersion.DownloadUrl)
		if err != nil {
			raven.CaptureError(err, map[string]string{"error": "Error downloading forge file: '" + modpack.Technic.ForgeVersion.DownloadUrl + "'"})
			conn.Log("Error downloading forge file: " + err.Error() + "\n" + modpack.Technic.ForgeVersion.DownloadUrl)
			return
		}
		defer resp.Body.Close()

		f, err := zipWriter.Create("bin/modpack.jar")
		if err != nil {
			raven.CaptureError(err, nil)
			conn.Log("Error while creating zip file content: " + err.Error())
			return
		}

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			raven.CaptureError(err, nil)
			conn.Log("Error white writing content to zip file: " + err.Error())
			return
		}
	}
	conn.Write(donePackingPartName, minecraftForge)

	if ch != nil {
		info := types.OutputInfo{
			File:             outputFile,
			Name:             minecraftForge,
			Id:               "forge",
			Version:          version,
			MinecraftVersion: modpack.MinecraftVersion,
			Description:      "The core of everything modded minecraft.",
			Url:              "http://www.minecraftforge.net/",
			Author:           "LexManos, cpw",
			ProgressKey:      minecraftForge,
		}
		*ch <- &info
	}
}

func packAdditionalFolder(modpack types.Modpack, folderPath string, outputDirectory string, conn types.WebsocketConnection, ch *chan *types.OutputInfo) {
	conn.Write(packingPartName, folderPath)
	inputFolder := path.Join(modpack.InputDirectory, folderPath)
	inputFolderInfo, _ := os.Stat(inputFolder)
	s := types.SafeNormalizeString(modpack.Name + "-" + inputFolderInfo.Name())
	outputDirectory = path.Join(outputDirectory, "mods", s)
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, s+"-"+modpack.GetVersionString()+".zip")
	zipfile, err := os.Create(outputFile)
	if err != nil {
		raven.CaptureError(err, nil)
		conn.Log("Error when creating zip file: " + err.Error() + "\n" + outputFile)
		return
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	packFolder(zipWriter, inputFolder, folderPath, conn)

	conn.Write(donePackingPartName, folderPath)

	if ch != nil {
		info := types.OutputInfo{
			File:             outputFile,
			Name:             modpack.Name + "-" + inputFolderInfo.Name(),
			Id:               s,
			Version:          modpack.Version,
			MinecraftVersion: modpack.MinecraftVersion,
			ProgressKey:      folderPath,
		}
		*ch <- &info
	}
}

func packFolder(zipWriter *zip.Writer, folder string, parent string, conn types.WebsocketConnection) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			packFolder(zipWriter, path.Join(folder, file.Name()), path.Join(parent, file.Name()), conn)
		} else {
			fileEntry := path.Join(parent, file.Name())

			fileReader, err := os.Open(path.Join(folder, file.Name()))
			defer fileReader.Close()

			if err != nil {
				raven.CaptureError(err, nil)
				conn.Log(err.Error())
				return
			}

			f, err := zipWriter.Create(fileEntry)
			if err != nil {
				raven.CaptureError(err, nil)
				conn.Log("Error while creating zip file content: " + err.Error() + "\n" + folder + "/" + file.Name())
				return
			}

			_, err = io.Copy(f, fileReader)
			if err != nil {
				raven.CaptureError(err, nil)
				conn.Log("Error white writing content to zip file: " + err.Error() + "\n" + folder + "/" + file.Name())
				return
			}
		}
	}
}

func packMod(mod *types.Mod, conn types.WebsocketConnection, outputDirectory string, ch *chan *types.OutputInfo) {
	if mod.Md5 == "" {
		fmt.Println("Calculating md5 of file " + mod.Filename)
		md5, _ := helpers.ComputeMd5(mod.Filename)
		mod.Md5 = hex.EncodeToString(md5)
	}
	conn.Write(packingPartName, mod.Filename)
	outputDirectory = path.Join(outputDirectory, "mods", mod.ModId)
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, mod.GetVersionString()+".zip")
	zipfile, err := os.Create(outputFile)
	if err != nil {
		raven.CaptureError(err, nil)
		conn.Log("Error when creating zip file: " + err.Error() + "\n" + outputFile)
		return
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	fileInfo, err := os.Stat(mod.Filename)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Println(err)
		conn.Error(err.Error())
	}
	file, err := os.Open(mod.Filename)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Println(err)
		conn.Error(err.Error())
	}
	defer file.Close()

	zipName := path.Join("mods", fileInfo.Name())

	f, err := zipWriter.Create(zipName)
	if err != nil {
		raven.CaptureError(err, nil)
		conn.Log("Error while creating zip file content: " + err.Error() + "\n" + outputFile)
		return
	}

	_, err = io.Copy(f, file)
	if err != nil {
		raven.CaptureError(err, nil)
		conn.Log("Error white writing content to zip file: " + err.Error() + "\n" + outputFile)
	}

	conn.Write(donePackingPartName, mod.Filename)

	if ch != nil {
		*ch <- GenerateOutputInfo(mod, outputFile)
	}
}

func GenerateOutputInfo(mod *types.Mod, outputFile string) *types.OutputInfo {
	info := types.OutputInfo{
		File:             outputFile,
		Name:             mod.Name,
		Id:               types.SafeNormalizeString(mod.ModId),
		Version:          mod.Version,
		MinecraftVersion: mod.MinecraftVersion,
		Description:      mod.Description,
		Author:           mod.Authors,
		ProgressKey:      mod.Filename,
		IsOnSolder:       mod.IsOnSolder,
		Permissions:      mod.Permission,
	}
	u := mod.Url
	if len(u) > 0 && strings.Index(u, "http") != 0 {
		u = "http://" + u
	}
	info.Url = u
	return &info
}
