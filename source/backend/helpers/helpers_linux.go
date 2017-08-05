package helpers

import (
	"log"
	"os/exec"
)

// Opens the given webpage in the default browser
func OpenWebPage(url string) {

	cmd := exec.Command("xdg-open", url)
	err := cmd.Run()
	if err != nil {
		log.Panicln(err)
	}
}
