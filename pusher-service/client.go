package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	id       string
	hub      *Hub
	socket   *websocket.Conn
	outbound chan []byte
}

// constructor
func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

// write to the websocket
func (c *Client) Write() {
	var err error
	for {
		// mux throught the outbound
		select {
		case message, ok := <-c.outbound:
			if !ok {
				// write to the socket that the conecction
				// is closed
				err = c.socket.WriteMessage(
					websocket.CloseMessage,
					[]byte{}, // an empty byte slice
				)
				if err != nil {
					log.Println(err)
				}
				return // finish
			}
			// write the message in the outbound
			err = c.socket.WriteMessage(
				websocket.TextMessage,
				message,
			)
			if err != nil {
				log.Println(err)
			}
		}
	}

}
