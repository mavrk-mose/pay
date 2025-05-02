package nats

import (
	"github.com/mavrk-mose/pay/pkg/utils"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	conn   *nats.Conn
	logger utils.Logger
}

// add these from a streams config -> so its like redis streams
const (
	StreamName     = "PAYMENTS" // TODO: these should be dynamic: PAYMENTS, NOTIFICATIONS, USER, TRANSACTION, etc
	StreamSubjects = "PAYMENTS.*" // TODO: as well as these
)

func JetStreamInit() (nats.JetStreamContext, error) {
	url := fmt.Sprintf("nats://%s:%s", c.Nats.NatsHost, c.Nats.NatsPort)
	
	//TODO: in dev use default url -> in prod use config's url
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256)) //PublishAsyncMaxPending sets the maximum outstanding async publishes that can be inflight at one time.
	if err != nil {
		return nil, err
	}

	// Create a stream if it does not exist
	err = CreateStream(js)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func CreateStream(jetStream nats.JetStreamContext) error {
	stream, err := jetStream.StreamInfo(StreamName)

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


//TODO: publishing to JetStream
// const (
// 	SubjectNameReviewCreated = "REVIEWS.rateGiven"
// )
//
// func publishReviews(js nats.JetStreamContext) {
// 	reviews, err := getReviews()
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	for _, oneReview := range reviews {
		
// 		// create random message intervals to slow down
// 		r := rand.Intn(1500)
// 		time.Sleep(time.Duration(r) * time.Millisecond)

// 		reviewString, err := json.Marshal(oneReview)
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}

// 		// publish to REVIEWS.rateGiven subject
// 		_, err = js.Publish(SubjectNameReviewCreated, reviewString)
// 		if err != nil {
// 			log.Println(err)
// 		} else {
// 			log.Printf("Publisher  =>  Message:%s\n", oneReview.Text)
// 		}
// 	}
// }

// func getReviews() ([]models.Review, error) {
// 	rawReviews, _ := ioutil.ReadFile("./reviews.json")
// 	var reviewsObj []models.Review
// 	err := json.Unmarshal(rawReviews, &reviewsObj)

// 	return reviewsObj, err
// }


//TODO: subscribing to JetStream
// func consumeReviews(js nats.JetStreamContext) {
// 	_, err := js.Subscribe(SubjectNameReviewCreated, func(m *nats.Msg) {
// 		err := m.Ack()

// 		if err != nil {
// 			log.Println("Unable to Ack", err)
// 			return
// 		}

// 		var review models.Review
// 		err = json.Unmarshal(m.Data, &review)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		log.Printf("Consumer  =>  Subject: %s  -  ID:%s  -  Author: %s  -  Rating:%d\n", m.Subject, review.Id, review.Author, review.Rating)
		
// 		// send answer via JetStream using another subject if you need
// 		// js.Publish(config.SubjectNameReviewAnswered, []byte(review.Id))
// 	})

// 	if err != nil {
// 		log.Println("Subscribe failed")
// 		return
// 	}
// }
