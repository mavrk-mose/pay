package service

import "context"

type Notifier interface {
	Send(ctx context.Context, userID, title string, details map[string]string) error
}
