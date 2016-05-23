package comm

import (
	"net/http"
	"flag"
	"log"
)


var addr = flag.String("addr", ":8084", "http service address")
func Run() {
	go h.run()
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
