package nats

import (
	"encoding/json"
	"go.uber.org/zap"

	"github.com/nats-io/nats.go"
)

// Publisher is a generic interface for message producers
// T is the type of the message you expect to publish
// It publishes the message asynchronously
type Publisher interface {
	Publish(subject string, data interface{}) error
}

// Consumer is a generic interface for message handlers
// T is the type of the message you expect to receive
// It gets called with the unmarshalled message and the raw *nats.Msg
// So that you can ACK or process metadata if needed
type Consumer interface {
	Consume(subject string, handler func(interface{})) error
}

type Client struct {
	JS     nats.JetStreamContext
	Logger *zap.Logger
}

func NewNatsClient(js nats.JetStreamContext) *Client {
	return &Client{
		JS:     js,
		Logger: zap.Must(zap.NewProduction()),
	}
}

func (n *Client) Publish(subject string, data interface{}) error {
	//n.Logger.Info("Publishing to subject: %s with data: %+v", zap.Any("subject", subject), zap.Any("data", data))

	msg, err := json.Marshal(data)
	if err != nil {
		n.Logger.Error("Failed to marshal message: %v", zap.Error(err))
		return err
	}

	_, err = n.JS.PublishAsync(subject, msg)
	if err != nil {
		n.Logger.Error("Failed to publish message: %v", zap.Error(err))
	}
	return err
}

func (n *Client) Consume(subject string, handler func(interface{})) error {
	n.Logger.Info("Subscribing to subject: %s", zap.Any("subject", subject))

	_, err := n.JS.Subscribe(subject, func(m *nats.Msg) {
		go func(msg *nats.Msg) {
			var data interface{}
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				n.Logger.Error("Failed to unmarshal message on subject %s: %v", zap.Any("subject", subject), zap.Error(err))
				return
			}

			n.Logger.Info("Received message on subject %s: %+v", zap.Any("subject", subject), zap.Any("data", data))

			handler(msg.Data)

			if err := msg.Ack(); err != nil {
				n.Logger.Error("Failed to acknowledge message on subject %s: %v", zap.Any("subject", subject), zap.Error(err))
			}
		}(m)
	}, nats.Durable("worker"))

	if err != nil {
		n.Logger.Error("Failed to subscribe to subject %s: %v", zap.Any("subject", subject), zap.Error(err))
	}
	return err
}
