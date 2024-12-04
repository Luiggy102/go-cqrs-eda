package search

import "github.com/elastic/go-elasticsearch/v7"

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

// methost for the searchRepository interface
// Close()
// IndexFeed(ctx context.Context, feed models.Feed) error
// SearchFeed(ctx context.Context, query string) ([]models.Feed, error)
