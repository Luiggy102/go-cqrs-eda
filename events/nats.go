package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/Luiggy102/go-cqrs-eda/models"
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
func (n *NatsEventStore) Close() {
	// check first for != nil
	if n.conn != nil {
		n.conn.Close()
	}
	if n.feedCreatedSub != nil {
		n.feedCreatedSub.Unsubscribe()
	}
	close(n.feedCreatedChan)
}

func (n *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	msg := CreatedFeedMessage{
		ID:          feed.ID,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}
	// sent the msg data to the nats Connection (publish)
	data, err := n.encodeMessage(msg)
	if err != nil {
		return err
	}

	// publish
	err = n.conn.Publish(msg.Type(), data)
	if err != nil {
		return err
	}
	return nil
}

func (n *NatsEventStore) OnCreateFeed(ctx context.Context, f func(CreatedFeedMessage)) error {
	var err error
	msg := CreatedFeedMessage{}
	n.feedCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		err = n.decodeMessage(m.Data, &msg)
		// call the function with the msg
		f(msg)
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *NatsEventStore) SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	m := CreatedFeedMessage{}
	n.feedCreatedChan = make(chan CreatedFeedMessage, 64)
	ch := make(chan *nats.Msg, 64) // only for msg
	// chan subscribe
	var err error
	n.feedCreatedSub, err = n.conn.ChanSubscribe(m.Type(), ch)
	if err != nil {
		return nil, err
	}
	// go routine
	go func() {
		// wait for a new msg
		for {
			select {
			case msg := <-ch:
				// decode and send
				n.decodeMessage(msg.Data, &m)
				n.feedCreatedChan <- m
			}
		}

	}()
	return (<-chan CreatedFeedMessage)(n.feedCreatedChan), nil
}
