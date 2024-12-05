package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Luiggy102/go-cqrs-eda/events"
	"github.com/Luiggy102/go-cqrs-eda/models"
	"github.com/Luiggy102/go-cqrs-eda/repository"
	"github.com/Luiggy102/go-cqrs-eda/search"
)

// index the msg to elastic search
func onCreatedFeed(m events.CreatedFeedMessage) {
	f := models.Feed{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
	err := search.IndexFeed(context.Background(), f)
	if err != nil {
		log.Printf("failed to index feed: %v", err)
	}
}

func ListFeedsHandler(w http.ResponseWriter, r *http.Request) {
	// feeds from the db
	feeds, err := repository.ListFeeds(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(feeds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
