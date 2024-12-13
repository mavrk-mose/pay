package ports

type WalletService interface {
	UpdateBalance(userID string, amount float64) error
	GetBalance(userID string) (float64, error)
}
