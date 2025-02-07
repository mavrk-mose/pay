package services

import (
	. "github.com/mavrk-mose/pay/internal/model"
)

type LedgerService interface {
	RecordTransaction(transaction LedgerTransaction) error
	GetTransactionByID(transactionID string) (LedgerTransaction, error)
	ReconcileTransactions(settlementFile SettlementFile) (ReconciliationResult, error)
}

func (l *LedgerService) RecordTransaction() {
	txn, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	txn := &models.Transaction{
		ExternalRef: externalRef,
		Status:      models.TransactionPending,
		Details:     nil, // Optionally store extra data as JSON (string)
	}

	// Create ledger entries: one debit and one credit
	debitEntry := models.Transaction{
		WalletID:  payerWalletID,
		EntryType: models.Debit,
		Amount:    amount,
		Currency:  currency,
	}
	creditEntry := models.Transaction{
		WalletID:  payeeWalletID,
		EntryType: models.Credit,
		Amount:    amount,
		Currency:  currency,
	}
	entries := []models.Transaction{debitEntry, creditEntry}

	// Save the transaction and its entries atomically
	if err := s.TransactionRepo.CreateTransactionWithEntries(ctx, txn, entries); err != nil {
		return "", err
	}
}
