package model

import "github.com/google/uuid"

type RiskScore struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Score     float64   `json:"score" db:"score"`
	RiskLevel string    `json:"risk_level" db:"risk_level"`
	Details   string    `json:"details" db:"details"`
}
