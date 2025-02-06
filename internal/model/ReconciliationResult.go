package model

type ReconciliationResult struct {
	Discrepancies []LedgerTransaction
	Matched       int
	Unmatched     int
}
