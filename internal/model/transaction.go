package model

import "time"

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
	Details     *string           `db:"details" json:"details,omitempty"` // Optional extra details as JSON or text
	Currency    string            `db:"currency" json:"currency"`

	// Debit side (e.g., the wallet being debited)
	DebitWalletID int64   		  `db:"debit_wallet_id" json:"debit_wallet_id"`
	DebitAmount   float64         `db:"debit_amount" json:"debit_amount"`

	EntryType     EntryType       `db:"entry_type" json:"entry_type"`

	// Credit side (e.g., the wallet being credited)
	CreditWalletID int64          `db:"credit_wallet_id" json:"credit_wallet_id"`
	CreditAmount   float64        `db:"credit_amount" json:"credit_amount"`

	CreatedAt time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt time.Time           `db:"updated_at" json:"updated_at"`
}