package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Luiggy102/go-cqrs-eda/events"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	// process Config
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	hub := NewHub()

	// add nats connection
	natsUrl := fmt.Sprintf("nats://%s", cfg.NatsAddress)
	nats, _ := events.NewNats(natsUrl)

	// send the msg to the hub
	err = nats.OnCreateFeed(context.Background(), func(m events.CreatedFeedMessage) {
		hub.Broadcast(
			newCreatedFeedMessage(m.ID, m.Title, m.Description, m.CreatedAt), nil,
		)
	})
	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStrore(nats)
	defer nats.Close()

	go hub.Run()

	// start the router
	http.HandleFunc("/ws", hub.HandleWebSocket)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
