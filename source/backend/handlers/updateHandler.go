package handlers

import (
	ghc "github.com/zlepper/github-release-checker"
	"github.com/zlepper/go-modpack-packer/source/backend/internal"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Checks for new versions of modpack packer
func CheckForUpdates(conn types.WebsocketConnection) {
	release, err := ghc.GetLatestReleaseForPlatform("zlepper", "go-modpack-packer", internal.FilenameRegex, false)
	if err != nil {
		log.Println("Could not check for updates", err)
		return
	}

	newer, err := ghc.IsNewer(release, internal.Version)
	if err != nil {
		log.Println("Could not compare release versions", err)
		return
	}

	if !newer {
		log.Printf("Running latest version '%s'\n", internal.Version)
		return
	}

	conn.Write("new-version-available", "")
}

const updateProgressWebsocketKey = "update-progress"

func updateToLatestVersion(conn types.WebsocketConnection) {
	conn.Write(updateProgressWebsocketKey, "Preparing update data")

	thisExecFileName := os.Args[0]

	wd, err := os.Getwd()
	if err != nil {
		log.Println("Could not get working directory", err)
		conn.Error(err)
		return
	}

	conn.Write(updateProgressWebsocketKey, "Getting release")
	release, err := ghc.GetLatestReleaseForPlatform("zlepper", "go-modpack-packer", internal.FilenameRegex, false)
	if err != nil {
		log.Println("Could not get latest release versions", err)
		conn.Error(err)
		return
	}

	absFileName, err := filepath.Abs(release.Filename)
	if err != nil {
		log.Println("Could not get output absolute filename", err)
		conn.Error(err)
		return
	}

	conn.Write(updateProgressWebsocketKey, "Downloading release")
	resp, err := http.Get(release.DownloadUrl)
	if err != nil {
		log.Println("Count not start download of new version", err)
		conn.Error(err)
		return
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(release.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println("Could not open new output file", err)
		conn.Error(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println("Could not write/download update contents", err)
		conn.Error(err)
		return
	}

	resp.Body.Close()
	file.Close()

	conn.Write(updateProgressWebsocketKey, "Removing old binary")

	err = os.Rename(thisExecFileName, thisExecFileName+".old")
	if err != nil {
		log.Println("Count not rename old binary", err)
		conn.Error(err)
		return
	}

	err = os.Rename(release.Filename, absFileName)

	conn.Write(updateProgressWebsocketKey, "Stopping old backend and starting new")
	port := internal.GetPort()

	internal.OutstandingProcess.Add(1)
	internal.EchoInstance.Close()

	processArgs := []string{absFileName, "--port", port, "--autolaunch=false", "--sendreloadsignal=true"}

	p, err := os.StartProcess(absFileName, processArgs, &os.ProcAttr{
		Dir: wd,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	})

	if err != nil {
		log.Panicln("Could not start new process for new version", err)
	} else {
		log.Println("New process started sucessfully")
	}
	conn.Write(updateProgressWebsocketKey, "New backend started")
	conn.Close()

	// Wait for new process to die off before killing this process, as the output pipes are still from this process
	p.Wait()

	internal.OutstandingProcess.Done()
}
