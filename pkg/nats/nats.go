package nats

import (
	"github.com/mavrk-mose/pay/pkg/utils"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	conn   *nats.Conn
	logger utils.Logger
}

func NewNATSClient(c *config.Config) (*Nats, error) {
  url := fmt.Sprintf("nats://%s:%s", c.Nats.NatsHost, c.Nats.NatsPort)
  
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Nats{conn: conn}, nil
}

// Publish sends a message to the given subject
func (n *Nats) Publish(subject string, data []byte) error {
	n.logger.Infof("Publishing message to subject %s", subject)
	return n.conn.Publish(subject, data)
}

// Subscribe sets up a handler for a given subject
func (n *Nats) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	n.logger.Infof("Subscribing to subject %s", subject)
	return n.conn.Subscribe(subject, handler)
}

// Close gracefully closes the connection
func (n *Nats) Close() {
	n.logger.Infof("Closing NATS connection")
	n.conn.Close()
}
