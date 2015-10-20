package kaca

import (
	"log"
	"net/http"
)

func ServeWs(addr string) {
	go disp.run()
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
