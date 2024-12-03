package events

import (
	"context"

	"github.com/Luiggy102/go-cqrs-eda/models"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error
	SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error)
	OnCreateFeed(ctx context.Context, f func(CreatedFeedMessage)) error
}

var eventStore EventStore

// set the abstract implementation
func SetEventStrore(store EventStore) {
	eventStore = store
}

func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, feed)
}
func SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	return eventStore.SubscribeCreatedFeed(ctx)
}
func OnCreatedFeed(ctx context.Context, f func(CreatedFeedMessage)) error {
	return eventStore.OnCreateFeed(ctx, f)
}
