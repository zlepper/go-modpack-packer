package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zlepper/go-modpack-packer/source/backend/solder"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/upload"
	wc "github.com/zlepper/go-websocket-connection"
	"runtime/debug"
	"log"
)

func HandleMessage(conn wc.WebsocketConnection, messageType int, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(fmt.Sprint(r) + "\n" + string(debug.Stack()))
			conn.Log(fmt.Sprint(r) + "\n" + string(debug.Stack()))
			conn.Write("notification", "Things just went skyhigh. Check the logs!!")
		}
	}()
	if messageType == websocket.TextMessage {
		var m wc.Message
		json.Unmarshal(message, &m)
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
		}
	}
}
