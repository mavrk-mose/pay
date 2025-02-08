package services

import (
	. "github.com/mavrk-mose/pay/internal/model"
	"time"
)

type ReconciliationService interface {
	ReconcileLedgerWithSettlement(ledger LedgerService, settlementFile SettlementFile) (ReconciliationResult, error)
	FlagPayment(payment PaymentIntent) error
	AnalyzeTransactions(ledger LedgerService) ([]Transaction, error)
	GenerateFraudReport(period time.Duration) (Report, error)
	EscalatePayment(payment PaymentIntent) error
	ReversePayment(payment PaymentIntent) error
	EvaluateUserRisk(userID int64) (RiskScore, error)
}
