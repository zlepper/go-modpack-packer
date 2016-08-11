package handlers

import (
	"archive/zip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/zlepper/go-modpack-packer/source/backend/db"
	"github.com/zlepper/go-modpack-packer/source/backend/helpers"
	"github.com/zlepper/go-modpack-packer/source/backend/solder"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/upload"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"github.com/zlepper/go-websocket-connection"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const donePackingPartName string = "done-packing-part"
const packingPartName string = "packing-part"

func build(conn websocket.WebsocketConnection, data interface{}) {
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

func continueRunning(conn websocket.WebsocketConnection, data interface{}) {
	dict := data.(map[string]interface{})

	var uploadInfo uploadWaiting
	err := mapstructure.Decode(dict, &uploadInfo)
	if err != nil {
		log.Panic(err)
	}

	solderclient, buildId := updateSolder(uploadInfo.Modpack, conn)
	conn.Write(solderCurrentlyDoingEvent, "BUILD.SOLDER.UPDATING_MODS")
	for _, info := range uploadInfo.Infos {
		go addInfoToSolder(info, buildId, conn, solderclient)
	}
}

func buildModpack(modpack types.Modpack, mods []*types.Mod, conn websocket.WebsocketConnection) {
	// Create output directory
	outputDirectory := path.Join(modpack.OutputDirectory, modpack.Name)
	os.MkdirAll(outputDirectory, os.ModePerm)
	var total int

	var ch chan *types.OutputInfo
	ch = make(chan *types.OutputInfo)

	startTime := time.Now()
	// Handle forge
	if modpack.Technic.CreateForgeZip {
		total++
		go packForgeFolder(modpack, conn, outputDirectory, &ch)
	}

	// Handle any additional folder
	for _, folder := range modpack.AdditionalFolders {
		if folder.Include {
			total++
			go packAdditionalFolder(modpack, folder.Name, outputDirectory, conn, &ch)
		}
	}

	infos := make([]*types.OutputInfo, 0)
	// Handle mods
	for _, mod := range mods {
		mod.NormalizeAll()
		// If the mod already is on solder, then we should likely skip it
		// however the user can override this. If they do we should still pack all files
		if !modpack.Technic.RepackAllMods && mod.IsOnSolder {
			conn.Write(packingPartName, mod.Filename)
			infos = append(infos, GenerateOutputInfo(mod, ""))
			total++
			conn.Write(donePackingPartName, mod.Filename)
			continue
		}
		go packMod(mod, conn, outputDirectory, &ch)
		total++
	}
	conn.Write("total-to-pack", total)

	// Save the mods to the database
	d := db.GetModsDb()
	for _, mod := range mods {
		mod.IsOnSolder = true
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
	}

	if modpack.Solder.Use {
		solderclient, buildId := updateSolder(modpack, conn)

		conn.Write(solderCurrentlyDoingEvent, "BUILD.SOLDER.UPDATING_MODS")
		for _, info := range infos {
			go addInfoToSolder(info, buildId, conn, solderclient)
		}
	}

}

const solderCurrentlyDoingEvent string = "solder-currently-doing"

func updateSolder(modpack types.Modpack, conn websocket.WebsocketConnection) (*solder.SolderClient, string) {
	// Create solder client
	conn.Write(solderCurrentlyDoingEvent, "BUILD.SOLDER.LOGGING_IN")
	solderclient := solder.NewSolderClient(modpack.Solder.Url)
	loginSuccessful := solderclient.Login(modpack.Solder.Username, modpack.Solder.Password)
	if !loginSuccessful {
		log.Panic("Could not login to solder with the supplied credentials.")
	}

	conn.Write(solderCurrentlyDoingEvent, "BUILD.SOLDER.CHECKING_MODPACK_EXISTENCE")
	var modpackId string
	if solderclient.IsPackOnline(&modpack) {
		conn.Write(solderCurrentlyDoingEvent, "BUILD.SOLDER.GETTING_MODPACK")
		modpackId = solderclient.GetModpackId(modpack.GetSlug())
	} else {
		conn.Write(solderCurrentlyDoingEvent, "BUILD.SOLDER.CREATING_MODPACK")
		modpackId = solderclient.CreatePack(modpack.Name, modpack.GetSlug())
	}

	var buildId string
	conn.Write(solderCurrentlyDoingEvent, "BUILD.SOLDER.CHECKING_BUILD_EXISTENCE")
	if solderclient.IsBuildOnline(&modpack) {
		conn.Write(solderCurrentlyDoingEvent, "BUILD.SOLDER.GETTING_BUILD")
		buildId = solderclient.GetBuildId(&modpack)
	} else {
		conn.Write(solderCurrentlyDoingEvent, "BUILD.SOLDER.CREATING_BUILD")
		buildId = solderclient.CreateBuild(&modpack, modpackId)
	}
	return solderclient, buildId
}

func addInfoToSolder(info *types.OutputInfo, buildId string, conn websocket.WebsocketConnection, solderclient *solder.SolderClient) {
	conn.Write("updating-solder", info.ProgressKey)
	var modid string
	modid = solderclient.GetModId(info.Id)
	if modid == "" {
		log.Println("Could not get mod id. Adding mod to solder " + info.Id)
		modid = solderclient.AddMod(info)
	}
	if modid == "" {
		log.Println("Something went wrong wehn adding a mod to solder.")
		log.Printf("%v\n", *info)
		log.Printf("Application version: %s\n", os.Args[2])
		log.Panic("Error. See above lines")
	}

	if info.File != "" {
		md5, err := helpers.ComputeMd5(info.File)
		if !solderclient.IsModversionOnline(info) {
			if err != nil {
				log.Panic(err)
			}
			solderclient.AddModVersion(modid, hex.EncodeToString(md5), info.GenerateOnlineVersion())
		} else {
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

func packForgeFolder(modpack types.Modpack, conn websocket.WebsocketConnection, outputDirectory string, ch *chan *types.OutputInfo) {
	const minecraftForge string = "Minecraft Forge"
	outputDirectory = path.Join(outputDirectory, "mods", "forge")
	os.MkdirAll(outputDirectory, os.ModePerm)
	version := fmt.Sprintf("%v", modpack.Technic.ForgeVersion.Build)
	outputFile := path.Join(outputDirectory, "forge-"+modpack.MinecraftVersion+"-"+version+".zip")
	conn.Write(packingPartName, minecraftForge)

	zipfile, err := os.Create(outputFile)
	if err != nil {
		conn.Log("Error when creating zip file: " + err.Error() + "\n" + outputFile)
		return
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	resp, err := http.Get(modpack.Technic.ForgeVersion.DownloadUrl)
	if err != nil {
		conn.Log("Error downloading forge file: " + err.Error() + "\n" + modpack.Technic.ForgeVersion.DownloadUrl)
		return
	}
	defer resp.Body.Close()

	f, err := zipWriter.Create("bin/modpack.jar")
	if err != nil {
		conn.Log("Error while creating zip file content: " + err.Error())
		return
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		conn.Log("Error white writing content to zip file: " + err.Error())
		return
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

func packAdditionalFolder(modpack types.Modpack, folderPath string, outputDirectory string, conn websocket.WebsocketConnection, ch *chan *types.OutputInfo) {
	conn.Write(packingPartName, folderPath)
	inputFolder := path.Join(modpack.InputDirectory, folderPath)
	inputFolderInfo, _ := os.Stat(inputFolder)
	s := safeNormalizeString(modpack.Name + "-" + inputFolderInfo.Name())
	outputDirectory = path.Join(outputDirectory, "mods", s)
	os.MkdirAll(outputDirectory, os.ModePerm)
	outputFile := path.Join(outputDirectory, s+"-"+modpack.GetVersionString()+".zip")
	zipfile, err := os.Create(outputFile)
	if err != nil {
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

func packFolder(zipWriter *zip.Writer, folder string, parent string, conn websocket.WebsocketConnection) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			packFolder(zipWriter, path.Join(folder, file.Name()), path.Join(parent, file.Name()), conn)
		} else {
			fileEntry := path.Join(parent, file.Name())

			fileReader, err := os.Open(path.Join(folder, file.Name()))
			defer fileReader.Close()

			if err != nil {
				conn.Log(err.Error())
				return
			}

			f, err := zipWriter.Create(fileEntry)
			if err != nil {
				conn.Log("Error while creating zip file content: " + err.Error() + "\n" + folder + "/" + file.Name())
				return
			}

			_, err = io.Copy(f, fileReader)
			if err != nil {
				conn.Log("Error white writing content to zip file: " + err.Error() + "\n" + folder + "/" + file.Name())
				return
			}
		}
	}
}

func packMod(mod *types.Mod, conn websocket.WebsocketConnection, outputDirectory string, ch *chan *types.OutputInfo) {
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
		conn.Log("Error when creating zip file: " + err.Error() + "\n" + outputFile)
		return
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	fileInfo, err := os.Stat(mod.Filename)
	if err != nil {
		log.Println(err)
		conn.Error(err.Error())
	}
	file, err := os.Open(mod.Filename)
	if err != nil {
		log.Println(err)
		conn.Error(err.Error())
	}
	defer file.Close()

	zipName := path.Join("mods", fileInfo.Name())

	f, err := zipWriter.Create(zipName)
	if err != nil {
		conn.Log("Error while creating zip file content: " + err.Error() + "\n" + outputFile)
		return
	}

	_, err = io.Copy(f, file)
	if err != nil {
		conn.Log("Error white writing content to zip file: " + err.Error() + "\n" + outputFile)
	}

	conn.Write(donePackingPartName, mod.Filename)

	if ch != nil {
		*ch <- GenerateOutputInfo(mod, outputFile)
	}
}

func safeNormalizeString(s string) string {
	s = strings.Replace(strings.ToLower(s), " ", "-", -1)
	return strings.Replace(s, ".", "", -1)
}

func GenerateOutputInfo(mod *types.Mod, outputFile string) *types.OutputInfo {
	info := types.OutputInfo{
		File:             outputFile,
		Name:             mod.Name,
		Id:               safeNormalizeString(mod.ModId),
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
