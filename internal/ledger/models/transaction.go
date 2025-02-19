package models

import (
	"github.com/google/uuid"
	"time"
)

type TransactionType string

const (
	TransactionWithdrawal TransactionType = "withdrawal"
	TransactionDeposit    TransactionType = "deposit"
	TransactionTransfer   TransactionType = "transfer"
	TransactionCharge     TransactionType = "charge"
)

type TransactionStatus string

const (
	TransactionPending   TransactionStatus = "pending"
	TransactionConfirmed TransactionStatus = "confirmed"
	TransactionFailed    TransactionStatus = "failed"
)

type EntryType string

const (
	Debit  EntryType = "debit"
	Credit EntryType = "credit"
)

type Transaction struct {
	ID          uuid.UUID         `db:"id" json:"id"`
	ExternalRef string            `db:"external_ref" json:"external_ref"` // Unique external reference to track the transaction
	Type        TransactionType   `db:"type" json:"type"`                 // e.g. transfer, deposit, withdrawal, charge
	Status      TransactionStatus `db:"status" json:"status"`             // pending, confirmed, or failed
	Details     string            `db:"details" json:"details,omitempty"` // Optional extra details as JSON or text
	Currency    string            `db:"currency" json:"currency"`

	DebitWalletID int64   `db:"debit_wallet_id" json:"debit_wallet_id"`
	DebitAmount   float64 `db:"debit_amount" json:"debit_amount"`

	EntryType EntryType `db:"entry_type" json:"entry_type"`

	CreditWalletID int64   `db:"credit_wallet_id" json:"credit_wallet_id"`
	CreditAmount   float64 `db:"credit_amount" json:"credit_amount"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Report struct {
	ReportID      uuid.UUID `json:"report_id" db:"report_id"`
	GeneratedAt   time.Time `json:"generated_at" db:"generated_at"`
	Period        string    `json:"period" db:"period"`
	TotalCases    int       `json:"total_cases" db:"total_cases"`
	TotalAmount   float64   `json:"total_amount" db:"total_amount"`
	HighRiskUsers []string  `json:"high_risk_users" db:"high_risk_users"`
	Details       string    `json:"details" db:"details"`
}

type SettlementFile struct {
	FileID        string
	BankAccountID string
	Transactions  []Transaction
}

type RiskScore struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Score     float64   `json:"score" db:"score"`
	RiskLevel string    `json:"risk_level" db:"risk_level"`
	Details   string    `json:"details" db:"details"`
}
