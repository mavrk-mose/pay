package service

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/mavrk-mose/pay/config"
	repo "github.com/mavrk-mose/pay/internal/notification/repository"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
	"google.golang.org/api/option"
)

type PushNotifier struct {
	repo     repo.NotificationRepo
	userRepo repository.UserRepository
	logger   utils.Logger
	firebase config.Firebase
}

func NewPushNotifier(
	repo repo.NotificationRepo,
	userRepo repository.UserRepository,
	cfg *config.Config,
) *PushNotifier {
	return &PushNotifier{
		repo:     repo,
		userRepo: userRepo,
		firebase: cfg.Firebase,
	}
}

func (n *PushNotifier) Send(ctx context.Context, user models.User, templateID string, details map[string]string) error {
	if user.DeviceToken == "" {
		n.logger.Warnf("No device token found for user %s", user.ID.String())
		return fmt.Errorf("no device token registered for user")
	}

	credJSON, err := json.Marshal(n.firebase)
	if err != nil {
		return fmt.Errorf("failed to marshal Firebase credentials: %v", err)
	}

	opt := option.WithCredentialsJSON(credJSON)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		n.logger.Errorf("Failed to initialize Firebase app: %v", err)
		return fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		n.logger.Errorf("Failed to get messaging client: %v", err)
		return fmt.Errorf("failed to get messaging client: %v", err)
	}

	template, err := n.repo.GetTemplate(ctx, templateID)
	if err != nil {
		n.logger.Errorf("Failed to get template %s: %v", templateID, err)
		return fmt.Errorf("failed to get template: %w", err)
	}

	message := utils.ReplaceTemplatePlaceholders(template.Message, details)

	n.logger.Debugf("Processed template message: %s", message)

	_, err = client.Send(ctx, &messaging.Message{
		Token: user.DeviceToken,
		Notification: &messaging.Notification{
			Title: template.Title,
			Body:  message,
		},
	})
	if err != nil {
		n.logger.Errorf("Failed to send push notification: %v", err)
		return fmt.Errorf("failed to send push notification: %v", err)
	}

	return nil
}

//
//import (
//"context"
//"log"
//"sync"
//"time"
//
//"firebase.google.com/go/messaging"
//)
//
//// Buffer batches all incoming push messages and send them periodically.
//type Buffer struct {
//	fcmClient *messaging.Client
//
//	dispatchInterval time.Duration
//	batchCh          chan *messaging.Message
//	wg               sync.WaitGroup
//}
//
//func (b *Buffer) SendPush(msg *messaging.Message) {
//	b.batchCh <- msg
//}
//
//func (b *Buffer) sender() {
//	defer b.wg.Done()
//
//	// set your interval
//	t := time.NewTicker(b.dispatchInterval)
//
//	// we can send up to 500 messages per call to Firebase
//	messages := make([]*messaging.Message, 0, 500)
//
//	defer func() {
//		t.Stop()
//
//		// send all buffered messages before quit
//		b.sendMessages(messages)
//
//		log.Println("batch sender finished")
//	}()
//
//	for {
//		select {
//		case m, ok := <-b.batchCh:
//			if !ok {
//				return
//			}
//
//			messages = append(messages, m)
//		case <-t.C:
//			b.sendMessages(messages)
//			messages = messages[:0]
//		}
//	}
//}
//
//func (b *Buffer) Run() {
//	b.wg.Add(1)
//	go b.sender()
//}
//
//func (b *Buffer) Stop() {
//	close(b.batchCh)
//	b.wg.Wait()
//}
//
//func (b *Buffer) sendMessages(messages []*messaging.Message) {
//	if len(messages) == 0 {
//		return
//	}
//
//	batchResp, err := b.fcmClient.SendAll(context.TODO(), messages)
//
//	log.Printf("batch response: %+v, err: %s \n", batchResp, err)
//}

// TODO: for push notification require device id & platform in the http handler
// deviceID := ctx.GetHeader("device-id")
// if deviceID == "" {
// 	response.BadRequestError(ctx, "device-id missing")
// 	return
// }
// platform := ctx.GetHeader("platform") -- do android only for now
