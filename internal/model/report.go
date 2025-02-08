package model

import (
	"github.com/google/uuid"
	"time"
)

type ReconciliationResult struct {
	Discrepancies []Transaction
	Matched       int
	Unmatched     int
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
