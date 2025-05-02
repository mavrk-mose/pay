package nats

import "fmt"

type Stream string
type Subject string

const (
	Payment       Stream = "PAYMENT"
	Ledger        Stream = "LEDGER"
	Wallet        Stream = "WALLET"
	Notifications Stream = "NOTIFICATIONS"
)

const (
	PaymentProcessed  Subject = "payment.processed"
	PaymentExecuted   Subject = "payment.executed"
	WalletFunded      Subject = "wallet.fund"
	UserCreated       Subject = "user.created"
	UserUpdated       Subject = "user.updated"
	UserDeleted       Subject = "user.deleted"
	UserWalletCreated Subject = "user.wallet.created"
	UserWalletUpdated Subject = "user.wallet.updated"
	UserWalletDeleted Subject = "user.wallet.deleted"
	UserNotification  Subject = "user.notification"
)

var schema = map[Stream][]Subject{
	Payment: {
		PaymentProcessed,
		PaymentExecuted,
	},
	Wallet: {
		WalletFunded,
	},
	Ledger: {},
	Notifications: {
		UserCreated,
		UserUpdated,
		UserDeleted,
		UserWalletCreated,
		UserWalletUpdated,
		UserWalletDeleted,
		UserNotification,
	},
}

func ValidateSubject(stream Stream, subject Subject) error {
	subjects, ok := schema[stream]
	if !ok {
		return fmt.Errorf("stream %q is not defined in schema", stream)
	}

	for _, s := range subjects {
		if s == subject {
			return nil
		}
	}
	return fmt.Errorf("subject %q is not valid under stream %q", subject, stream)
}
