package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Luiggy102/go-cqrs-eda/events"
	"github.com/Luiggy102/go-cqrs-eda/models"
	"github.com/Luiggy102/go-cqrs-eda/repository"
	"github.com/segmentio/ksuid"
)

type createFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req createFeedRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// newId
	id, err := ksuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// new feed
	feed := models.Feed{
		ID:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now().UTC(),
	}

	// add to the db
	err = repository.InsertFeed(r.Context(), &feed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// notify recently created feed through nats
	// publish event
	err = events.PublishCreatedFeed(r.Context(), &feed)
	if err != nil {
		log.Printf("Failed to publish created feed event: %v\n", err)
	}

	// send response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(feed)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
