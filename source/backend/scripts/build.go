// Builds the releases

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type configuration struct {
	OS        string
	Arch      string
	Extension string
}

var (
	configurations []configuration = []configuration{
		{
			OS:        "windows",
			Arch:      "386",
			Extension: "windows-x32.exe",
		},
		{
			OS:        "windows",
			Arch:      "amd64",
			Extension: "windows-x64.exe",
		},
		{
			OS:        "darwin",
			Arch:      "386",
			Extension: "osx-x32",
		},
		{
			OS:        "darwin",
			Arch:      "amd64",
			Extension: "osx-x64",
		},
		{
			OS:        "linux",
			Arch:      "386",
			Extension: "linux-x32",
		},
		{
			OS:        "linux",
			Arch:      "amd64",
			Extension: "linux-x64",
		},
	}
)

func main() {
	log.Println("Starting build")
	goBinary, err := exec.LookPath("go")
	if err != nil {
		log.Panicln(err)
	}
	for _, conf := range configurations {
		log.Printf("building binary for '%s'\n", conf.Extension)
		cmd := exec.Cmd{
			Path: goBinary,
			Args: []string{
				goBinary,
				"build",
				"-o",
				fmt.Sprintf("build/modpack-packer-%s", conf.Extension),
				"github.com/zlepper/go-modpack-packer/source/backend",
			},
			Env: append(
				os.Environ(),
				fmt.Sprintf("GOOS=%s", conf.OS),
				fmt.Sprintf("GOARCH=%s", conf.Arch),
			),
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("build args", cmd.Args)
			log.Println("Error when building", err, "\n", string(output))
			return
		}
		log.Printf("Successfully build binary for '%s'\n", conf.Extension)
	}
	log.Println("Finished building all configurations.")
}
