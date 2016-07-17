package main

import (
	"github.com/zlepper/go-modpack-packer/source/backend/handlers"
	"github.com/zlepper/go-websocket-connection"
	"log"
	"os"
	"path"
)

func main() {
	logfilePath := path.Join(os.Args[1], "log.log")
	file, err := os.OpenFile(logfilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	log.SetOutput(file)
	websocket.Run(handlers.HandleMessage)
}
