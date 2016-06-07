package main

import (
	"github.com/zlepper/go-modpack-packer/source/backend/comm"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	log.SetOutput(file)
	comm.Run()
}
