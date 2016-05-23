package handlers

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
)

type Message struct {
	Action string `json:"action"`
	Data   interface{} `json:"data"`
}

type websocketConnection interface {
	Log(message string)
	Write(data interface{})
}

func HandleMessage(conn websocketConnection, messageType int, message []byte) {
	log.Println("TEST")
	if messageType == websocket.TextMessage {
		var m Message
		json.Unmarshal(message, &m)
		switch m.Action {
		case "find-additional-folders": {
			findAdditionalFolders(conn, m)
		}
		}
	}
}

