package ledger

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/ledger/models"
	"sync"
	"time"
)

type TransactionRepo interface {
	RecordTransaction(ctx *gin.Context, payerWalletID, payeeWalletID int64, amount float64, currency string) (string, error)
	CreateTransactionWithEntries(ctx *gin.Context, txn *sqlx.Tx, entries []Transaction) error
	UpdateTransactionStatus(ctx *gin.Context, externalRef string, status TransactionStatus) error
	FetchTransactionsWithChecksum(db *sqlx.DB, date, provider string) (map[string]string, error)
}

type Repo struct {
	DB *sqlx.DB
}

func (r *Repo) RecordTransaction(ctx *gin.Context, payerWalletID, payeeWalletID int64, amount float64, currency string) (string, error) {
	txn, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	transactionID := uuid.New()

	debitEntry := Transaction{
		ID:            transactionID,
		ExternalRef:   uuid.New().String(),
		Type:          TransactionTransfer,
		Status:        TransactionPending,
		Currency:      currency,
		DebitWalletID: payerWalletID,
		Amount:        amount,
		EntryType:     Debit,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	creditEntry := Transaction{
		ID:             transactionID,
		ExternalRef:    uuid.New().String(),
		Type:           TransactionTransfer,
		Status:         TransactionPending,
		Currency:       currency,
		CreditWalletID: payeeWalletID,
		Amount:         amount,
		EntryType:      Credit,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	entries := []Transaction{debitEntry, creditEntry}

	if err := r.CreateTransactionWithEntries(ctx, txn, entries); err != nil {
		err := txn.Rollback()
		if err != nil {
			return "", err
		}
		return "", err
	}

	if err := txn.Commit(); err != nil {
		return "", err
	}

	return transactionID.String(), nil
}

func (r *Repo) CreateTransactionWithEntries(ctx *gin.Context, txn *sqlx.Tx, entries []Transaction) error {
	for i := range entries {
		entries[i].Checksum = GenerateChecksum(entries[i])

		_, err := txn.NamedExecContext(ctx, `
			INSERT INTO transaction (
				id, external_ref, type, status, details, currency, 
				debit_wallet_id, amount,
				credit_wallet_id, credit_amount, created_at, updated_at, checksum
			) VALUES (
				:id, :external_ref, :type, :status, :details, :currency, 
				:debit_wallet_id, :amount,
				:credit_wallet_id, :created_at, :updated_at, :checksum
			)`, entries[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) UpdateTransactionStatus(ctx *gin.Context, externalRef string, status TransactionStatus) error {
	query := `
		UPDATE transaction
		SET status = $1, updated_at = NOW() 
		WHERE external_ref = $2
	`
	res, err := r.DB.ExecContext(ctx, query, status, externalRef)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		return fmt.Errorf("no transaction found with external_ref %s", externalRef)
	}
	return nil
}

func FetchTransactionsWithChecksum(db *sqlx.DB, date, provider string) (map[string]string, error) {
	query := `SELECT id, checksum FROM transaction WHERE provider = $1 AND created_at >= NOW() - INTERVAL '1 DAY'`
	rows, err := db.Query(query, provider)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	dbTransactions := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for rows.Next() {
		var txn Transaction
		if err := rows.Scan(&txn.ID, &txn.Checksum); err != nil {
			panic(err)
		}
		wg.Add(1)
		go func(txn Transaction) {
			defer wg.Done()
			mu.Lock()
			dbTransactions[txn.ID.String()] = txn.Checksum
			mu.Unlock()
		}(txn)
	}
	wg.Wait()

	return dbTransactions, nil
}

func GenerateChecksum(txn Transaction) string {
	data := fmt.Sprintf("%s|%s|%s|%s|%f|%s|%s",
		txn.ExternalRef, txn.Type, txn.Status, txn.Currency,
		txn.Amount, txn.DebitWalletID, txn.CreditWalletID,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// TODO: add a transaction record for transfers
// // Insert the transfer record
// insertQuery := `
// INSERT INTO transfers (from_wallet_id, to_wallet_id, amount, currency, status, external_ref)
// VALUES (:from_wallet_id, :to_wallet_id, :amount, :currency, :status, :external_ref)
// `
// _, err = tx.NamedExecContext(ctx, insertQuery, transfer)
// if err != nil {
// tx.Rollback()
// r.logger.Errorf("Failed to create transfer record: %v", err)
// return fmt.Errorf("failed to create transfer record: %v", err)
// }

// TODO: add withdraw transaction
// debitQuery := `
// 		INSERT INTO transaction (transaction_id, wallet_id, entry_type, amount, currency)
// 		VALUES ($1, $2, $3, $4, $5)
// 	`
// 	_, err = tx.ExecContext(ctx, debitQuery, transactionID, walletID, "DEBIT", amount, currency)
// 	if err != nil {
// 		tx.Rollback()
// 		return "", err
// 	}

// creditQuery := `
// 		INSERT INTO transaction (transaction_id, wallet_id, entry_type, amount, currency)
// 		VALUES ($1, $2, $3, $4, $5)
// 	`
// 	_, err = tx.ExecContext(ctx, creditQuery, transactionID, walletID, "CREDIT", amount, currency)
// 	if err != nil {
// 		tx.Rollback()
// 		return "", err
// 	}
