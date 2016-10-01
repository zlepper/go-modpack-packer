package main

import (
	"github.com/getsentry/raven-go"
	"github.com/zlepper/go-modpack-packer/source/backend/handlers"
	"github.com/zlepper/go-websocket-connection"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	raven.SetDSN("https://68c8787167d940b1b2bd2e6a8308f242@app.getsentry.com/92966")

	raven.CapturePanic(func() {
		logfilePath := path.Join(os.Args[1], "log.log")
		file, err := os.OpenFile(logfilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		log.SetOutput(io.MultiWriter(file, os.Stdout))
		websocket.Run(handlers.HandleMessage)
	}, nil)
}
