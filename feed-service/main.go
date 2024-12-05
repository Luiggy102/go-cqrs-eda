package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Luiggy102/go-cqrs-eda/database"
	"github.com/Luiggy102/go-cqrs-eda/events"
	"github.com/Luiggy102/go-cqrs-eda/repository"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

// Config for start the (cqrs-command) feed service
type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/feeds", createFeedHandler).Methods(http.MethodPost)
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

	// add nats connection
	natsUrl := fmt.Sprintf("nats://%s", cfg.NatsAddress)
	nats, _ := events.NewNats(natsUrl)
	events.SetEventStrore(nats)
	defer nats.Close()

	// start the router
	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln(err)
	}
}
