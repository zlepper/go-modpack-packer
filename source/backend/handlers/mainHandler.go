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

func Write(conn websocketConnection, action string, data interface{}) {
	message := Message{
		Action:action,
		Data:data,
	}
	conn.Write(message)
}

func HandleMessage(conn websocketConnection, messageType int, message []byte) {
	log.Println("TEST")
	if messageType == websocket.TextMessage {
		var m Message
		json.Unmarshal(message, &m)
		switch m.Action {
		case "find-additional-folders": {
			findAdditionalFolders(conn, m.Data)
		}
		case "save-modpacks": {
			log.Printf("%v", m.Data)
			saveModpacks(conn, m.Data)
		}
		case "load-modpacks": {
			loadModpacks(conn)
		}
		}
	}
}

