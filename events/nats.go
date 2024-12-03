package events

import (
	"bytes"
	"encoding/gob"

	"github.com/nats-io/nats.go"
)

// event store
type NatsEventStore struct {
	conn            *nats.Conn
	feedCreatedSub  *nats.Subscription
	feedCreatedChan chan CreatedFeedMessage
}

// constructor
func NewNats(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{conn: conn}, nil
}

// encode/decode message
func (n *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (n *NatsEventStore) decodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	err := gob.NewDecoder(&b).Decode(m)
	if err != nil {
		return err
	}
	return nil
}

// satisfy the EventStore interface
// func (n *NatsEventStore) Close()                                                          {}
// func (n *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {}
// func (n *NatsEventStore) SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
// }
// func (n *NatsEventStore) OnCreatedFeed(ctx context.Context, f func(CreatedFeedMessage)) error {}
