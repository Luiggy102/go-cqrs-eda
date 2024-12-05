package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Luiggy102/go-cqrs-eda/database"
	"github.com/Luiggy102/go-cqrs-eda/events"
	"github.com/Luiggy102/go-cqrs-eda/repository"
	"github.com/Luiggy102/go-cqrs-eda/search"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

// Config for start the (cqrs-query) feed service
type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticsearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/feeds", ListFeedsHandler).Methods(http.MethodGet)
	router.HandleFunc("/search", ListFeedsHandler).Methods(http.MethodGet)
	return
}

func main() {
	// process Config
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// add pg database
	pgUrl := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB,
	)
	db, _ := database.NewPostgresRepo(pgUrl)
	repository.SetRepository(db)

	// add elasticsearch
	esUrl := fmt.Sprintf("http://%s", cfg.ElasticsearchAddress)
	es, _ := search.NewElastic(esUrl)
	search.SetRepo(es)
	defer es.Close()

	// add nats connection
	natsUrl := fmt.Sprintf("nats://%s", cfg.NatsAddress)
	nats, _ := events.NewNats(natsUrl)
	// listen and react the created feed
	err = nats.OnCreateFeed(context.Background(), onCreatedFeed)
	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStrore(nats)
	defer nats.Close()

	// start the router
	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln(err)
	}
}
