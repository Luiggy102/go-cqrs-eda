package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

// upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// constructor
func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

// handle websocket func (only upgrade the connection)
func (hub *Hub) HandleWebSocket(w http.ResponseWriter,
	r *http.Request) {
	// upgrade the http to websocket
	socker, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error upgrading connection",
			http.StatusInternalServerError)
		return
	}
	// upgrade succesful
	// new client
	c := NewClient(hub, socker)
	// register the new client
	hub.register <- c

	// the cliente begin to write to the websocket
	go c.Write()
}

// handle the ws connections
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register: // add the client on connection
			hub.onConnect(client)
		case client := <-hub.unregister: // unregister the client
			hub.onDisconnect(client)
		}
	}
}

// log the client connection
// add an id
// append to the hub's client slice
func (hub *Hub) onConnect(client *Client) {
	log.Println("New connection:", client.socket.RemoteAddr())

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	client.id = client.socket.RemoteAddr().String()
	hub.clients = append(hub.clients, client)
}

// log the client disconnection
// delete from hub's client slice
func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client disconnected:", client.socket.RemoteAddr())

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	var index int // index for the disconnected client
	for i, c := range hub.clients {
		if c.id == client.id {
			index = i
		}
	}
	// the the client from the slice
	hub.clients = append(hub.clients[:index], hub.clients[index+1:]...)
}

// send a message throught the ws connection
func (hub *Hub) Broadcast(message interface{}, ignore *Client) {
	// mashal de content message
	data, _ := json.Marshal(message)

	// send to the clients
	for _, c := range hub.clients {
		if c != ignore {
			c.outbound <- data // send the data
		}
	}

}
