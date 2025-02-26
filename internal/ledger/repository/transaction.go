package ledger

import (
	"database/sql"
	"fmt"
	. "github.com/mavrk-mose/pay/internal/ledger/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/pkg/utils"
)

type Repo struct {
	DB *sqlx.DB
}

func (r *Repo) RecordTransaction(ctx *gin.Context, payerWalletID, payeeWalletID int64, amount float64, currency string) (string, error) {
	txn, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	transactionID := uuid.New()

	// Create ledger entries: one debit and one credit
	debitEntry := Transaction{
		ID:            transactionID,
		ExternalRef:   uuid.New().String(),
		Type:          TransactionTransfer,
		Status:        TransactionPending,
		Currency:      currency,
		DebitWalletID: payerWalletID,
		DebitAmount:   amount,
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
		CreditAmount:   amount,
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
			INSERT INTO transactions (
				id, external_ref, type, status, details, currency, 
				debit_wallet_id, debit_amount, entry_type, 
				credit_wallet_id, credit_amount, created_at, updated_at, checksum
			) VALUES (
				:id, :external_ref, :type, :status, :details, :currency, 
				:debit_wallet_id, :debit_amount, :entry_type, 
				:credit_wallet_id, :credit_amount, :created_at, :updated_at, :checksum
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

func FetchTransactionsWithChecksum(db *sql.DB, date string) (map[string]string, error) {
	query := `SELECT id, checksum FROM transactions WHERE provider = $1 AND created_at >= NOW() - INTERVAL '1 DAY'`
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
			dbTransactions[txn.ID] = txn.Checksum
			mu.Unlock()
		}(txn)
	}
	wg.Wait()

	return dbTransactions, nil
}


func GenerateChecksum(txn Transaction) string {
	data := fmt.Sprintf("%s|%s|%s|%s|%f|%f|%s|%s",
		txn.ExternalRef, txn.Type, txn.Status, txn.Currency,
		txn.DebitAmount, txn.CreditAmount, txn.DebitWalletID, txn.CreditWalletID,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}