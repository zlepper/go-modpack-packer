package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/s3"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"runtime/debug"
)

func HandleMessage(conn types.WebsocketConnection, messageType int, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			conn.Log(fmt.Sprint(r) + "\n" + string(debug.Stack()))
		}
	}()
	if messageType == websocket.TextMessage {
		var m types.Message
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
				s3.GetAwsBuckets(conn, m.Data)
			}
		}
	}
}
