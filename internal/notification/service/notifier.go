package service

import (
	"context"
	"github.com/mavrk-mose/pay/internal/user/models"
)

type Notifier interface {
	Send(ctx context.Context, user models.User, templateID string, details map[string]string) error
}
