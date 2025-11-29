package main

import (
	"fmt"
	"net/http"

	"holohub-be/ws"
)

func main() {
	http.HandleFunc("/ws", ws.WsHandler)
	go ws.HandleMessages()
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
