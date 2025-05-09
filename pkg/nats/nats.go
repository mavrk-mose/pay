package nats

import (
	"fmt"
	"log"

	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/pkg/utils"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	conn   *nats.Conn
	logger utils.Logger
}

const (
	StreamName     = "PAYMENTS"
	StreamSubjects = "PAYMENTS.*"
)

func JetStreamInit(cfg *config.Config) (nats.JetStreamContext, error) {
	url := fmt.Sprintf("nats://%s:%s", cfg.Nats.Host, cfg.Nats.Port)

	nc, err := nats.Connect(url, nats.UserInfo(cfg.Nats.User, cfg.Nats.Password))
	if err != nil {
		return nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256)) //PublishAsyncMaxPending sets the maximum outstanding async publishes that can be inflight at one time.
	if err != nil {
		return nil, err
	}

	client := NewNatsClient(js)

	err = CreateStream(client.JS)
	if err != nil {
		return nil, err
	}

	return js, nil
}

// TODO: will need to update this when I have multiple streams each with multiple subjects
func CreateStream(jetStream nats.JetStreamContext) error {
	stream, err := jetStream.StreamInfo(StreamName)
	if err != nil {
		log.Printf("Stream %s not found: %v\n", StreamName, err)
	}

	// stream not found, create it
	if stream == nil {
		log.Printf("Creating stream: %s\n", StreamName)

		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{StreamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Close gracefully closes the connection
func (n *Nats) Close() {
	n.logger.Infof("Closing NATS connection ...")
	n.conn.Close()
}
