package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
	}

	cookie, err := req.Cookie("auth")
	if err != nil {
		log.Println("cookie is not founded :", err)
		return
	}

	value, err := url.QueryUnescape(cookie.Value)

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(value),
	}

	log.Println("created new client")
	r.join <- client
	defer func() {
		r.leave <- client
		log.Println("leave channel sent")
	}()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			log.Println("client join...")
		case client := <-r.leave:
			log.Println("client leave...")
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			log.Println("message received...")
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}
