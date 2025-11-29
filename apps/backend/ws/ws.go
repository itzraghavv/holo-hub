package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadCast = make(chan []byte)
var mtx = &sync.Mutex{}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading: ", err)
	}

	defer conn.Close()

	mtx.Lock()
	clients[conn] = true
	mtx.Unlock()

	for {
		_, rawMsg, err := conn.ReadMessage()
		if err != nil {
			mtx.Lock()
			delete(clients, conn)
			mtx.Unlock()
			break
		}

		var msg Message

		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Println("invalid json: ", err)
			continue
		}
		broadCast <- rawMsg
	}
}

func HandleMessages() {
	for {
		message := <-broadCast

		mtx.Lock()

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mtx.Unlock()
	}
}
