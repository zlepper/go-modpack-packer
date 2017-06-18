package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zlepper/go-modpack-packer/source/backend/solder"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/upload"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"log"
	"runtime/debug"
)

func HandleMessage(conn types.WebsocketConnection, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(fmt.Sprint(r) + "\n" + string(debug.Stack()))
			conn.Log(fmt.Sprint(r) + "\n" + string(debug.Stack()))
			conn.Write("notification", "The backend just exploded. Check the logs!!")
		}
	}()
	var m types.Message
	err := json.Unmarshal(message, &m)
	if err != nil {
		panic(err)
	}
	log.Println(m.Action)
	switch m.Action {
	case "find-additional-folders":
		{
			findAdditionalFolders(conn, m.Data)
		}
	case "save-modpacks":
		{
			saveModpacks(conn, m.Data)
		}
	case "load-modpacks":
		{
			loadModpacks(conn)
		}
	case "gather-information":
		{
			gatherInformation(conn, m.Data)
		}
	case "build":
		{
			build(conn, m.Data)
		}
	case "get-aws-buckets":
		{
			upload.GetAwsBuckets(conn, m.Data)
		}
	case "test-ftp":
		{
			upload.TestFtp(conn, m.Data)
		}
	case "test-solder":
		{
			solder.TestSolderConnection(conn, m.Data)
		}
	case "continue-running":
		{
			continueRunning(conn, m.Data)
		}
	case "check-permission-store":
		{
			CheckPermissionStore(conn, m.Data)
		}

	default:
		{
			log.Println("Unknown action", m.Action)
		}
	}
}
