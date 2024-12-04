package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

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
func (r *ElacticSearchRepository) Close() {
	//
}
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

func (r *ElacticSearchRepository) SearchFeed(ctx context.Context, query string) ([]models.Feed, error) {
	results := []models.Feed{}
	var buf bytes.Buffer // for save the json
	/* elasticsearch multi-match query format
	    {
		  "query": {
		    "multi_match" : {
		      "query":    "this is a test",
		      "fields": [ "subject", "message" ],
		           "fuzzines":         3,
		           "cutoff_frequency": 0.0001,
		    }
		  }
		}
	*/
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"title", "description"},
				"fuzzines":         3,
				"cutoff_frequency": 0.0001,
			},
		},
	}

	// encode into the buffer
	err := json.NewEncoder(&buf).Encode(searchQuery)
	if err != nil {
		return nil, err
	}

	// do the search
	res, err := r.client.Search(
		r.client.Search.WithIndex("feeds"),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithContext(ctx),
		r.client.Search.WithTrackTotalHits(true),
	)
	// check for result error
	if err != nil || res.IsError() {
		return nil, errors.New("elasticsearch error " + res.String())
	}

	// close the respose body
	defer func() {
		err := res.Body.Close()
		if err != nil {
			results = nil
		}
	}()

	// decode the json response
	var jsonRes map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&jsonRes)
	if err != nil {
		return nil, err
	}

	/* elasticsearch respose format
	    {
		    ...
		    "hits": {
		        ...
		        "hits": [
		            {
		                ...
		                "_source":{
		                    fields
		                    ...
		                }
		            }
		        ]
		    }

		} */
	// range the json response for hits []interface{}
	for _, hit := range jsonRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		f := models.Feed{}

		// find "_source"
		source := hit.(map[string]interface{})["_source"]

		// marshall "_source"
		b, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		// unmarshall at the feed model
		err = json.Unmarshal(b, &f)
		if err != nil {
			return nil, err
		}

		// append to the results
		results = append(results, f)
	}

	return results, nil
}
