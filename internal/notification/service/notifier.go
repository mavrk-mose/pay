package service

import "context"

type Notifier interface {
	Send(ctx context.Context, userID, templateID string, details map[string]string) error
}
