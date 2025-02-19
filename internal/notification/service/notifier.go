package service

type Notifier interface {
	Send(userID, title, message string) error
}
