package search

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/Luiggy102/go-cqrs-eda/models"
	"github.com/elastic/go-elasticsearch/v7"
)

type ElacticSearchRepository struct {
	client *elasticsearch.Client
}

func NewElastic(url string) (*ElacticSearchRepository, error) {
	// new es instance
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}
	return &ElacticSearchRepository{client: client}, nil
}

// methods for the searchRepository interface
// func (r *ElacticSearchRepository) Close(){}
func (r *ElacticSearchRepository) IndexFeed(ctx context.Context, feed models.Feed) error {
	// create a json bytes for the feed
	body, _ := json.Marshal(feed)

	// index the feed to the client
	_, err := r.client.Index(

		"feeds",               // index name
		bytes.NewReader(body), // reader with the data

		// options
		r.client.Index.WithDocumentID(feed.ID), // add an id to the document
		r.client.Index.WithContext(ctx),        // contex for debug
		r.client.Index.WithRefresh("wait_for"), // refresh the indices
	)

	if err != nil {
		return err
	}
	return nil
}

// func (r *ElacticSearchRepository) SearchFeed(ctx context.Context, query string) ([]models.Feed, error){}
