package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
  	"github.com/mavrk-mose/pay/pkg/utils"
)

// Publisher is a generic interface for message producers
// T is the type of the message you expect to publish
// It publishes the message asynchronously
type Publisher interface {
	Publish[T any](subject string, data T) error
}

// Consumer is a generic interface for message handlers
// T is the type of the message you expect to receive
// It gets called with the unmarshalled message and the raw *nats.Msg
// So that you can ACK or process metadata if needed
type Consumer interface {
	Consume[T any](subject string, handler func(T)) error
}

type NatsClient struct {
	JS 	nats.JetStreamContext
  	logger  utils.Logger
}

func NewNatsClient(js nats.JetStreamContext) *NatsClient {
	return &NatsClient{JS: js}
}

func (n *NatsClient) Publish[T any](subject string, data T) error {
	n.Logger.Infof("Publishing to subject: %s with data: %+v", subject, data)

	msg, err := json.Marshal(data)
	if err != nil {
		n.Logger.Errorf("Failed to marshal message: %v", err)
		return err
	}

	_, err = n.JS.PublishAsync(subject, msg)
	if err != nil {
		n.Logger.Errorf("Failed to publish message: %v", err)
	}
	return err
}

func (n *NatsClient) Consume[T any](subject string, handler func(T)) error {
	n.Logger.Infof("Subscribing to subject: %s", subject)

	_, err := n.JS.Subscribe(subject, func(m *nats.Msg) {
		go func(msg *nats.Msg) {
			var data T
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				n.Logger.Errorf("Failed to unmarshal message on subject %s: %v", subject, err)
				return
			}

			n.Logger.Infof("Received message on subject %s: %+v", subject, data)

			handler(data)

			if err := msg.Ack(); err != nil {
				n.Logger.Errorf("Failed to acknowledge message on subject %s: %v", subject, err)
			}
		}(m)
	}, nats.Durable("worker"))

	if err != nil {
		n.Logger.Errorf("Failed to subscribe to subject %s: %v", subject, err)
	}
	return err
}

