package types

type WebsocketConnection interface {
	Log(message string)
	WriteData(data interface{})
	Write(action string, data interface{})
	Error(data interface{})
}

type Message struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}
