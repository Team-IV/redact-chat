package main

import (
	"flag"
	"net/http"

	"golang.org/x/net/websocket"
)

// Message is going to hold...Messages.
type Message struct {
	Text string `json:"Text"`
}

// Hub is the middle man between client and server.
type Hub struct {
	clients          map[string]*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan    chan Message
}

// This func will handle client requests to the hub.
func handler(ws *websocket.Conn, h *Hub) {
go h.run()
h.addClientChan <- ws

for {

	var m Message
	err := websocket.JSON.Receive(ws, &m)
	if err != nil {

		h.broadcastChan <- Message{err.Error()}
		h.removeClient(ws)
		return
	}
	h.broadcastChan <- m
}
}



}

var (
	port = flag.String("port", "9000", "Port used for websocket connection")
)

// This function builds a custom server.
func server(port string) error {
	h := newHub()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		handler(ws, h)
	}))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux}
	return s.ListenAndServe()
}
