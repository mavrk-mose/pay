package notifiers

import (
	"context"
	"github.com/mavrk-mose/pay/internal/user/models"
)

//go:generate mockery --name=Notifier --output=./mocks --filename=notifier.go --with-expecter
type Notifier interface {
	Send(ctx context.Context, user models.User, templateID string, details map[string]string) error
}
