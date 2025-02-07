package services

type ReconciliationService interface {
	ReconcileLedgerWithSettlement(ledger LedgerService, settlementFile SettlementFile) (ReconciliationResult, error)
    FlagPayment(payment Payment) error
    AnalyzeTransactions(ledger LedgerService) ([]SuspiciousTransaction, error)
    GenerateFraudReport(period time.Duration) (FraudReport, error)
    EscalatePayment(payment Payment) error
    ReversePayment(payment Payment) error
	EvaluateUserRisk(userID int64) (RiskScore, error)
}