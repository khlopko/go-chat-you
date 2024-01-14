package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil) 
    if err != nil {
    }
    go func() {
        for {
            _, messageBytes, err := conn.ReadMessage()
            if err != nil {
                break
            }
            message := string(messageBytes)
            fmt.Println(message)
        }
    }()
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        indexHandler(w, r)
    })
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        wsHandler(w, r)
    })
    server := &http.Server{
        Addr: ":8080",
        ReadHeaderTimeout: time.Second * 60,
    }
    err := server.ListenAndServe()
    if err != nil {
        log.Fatal("Failed to start a server: ", err)
    }
}
