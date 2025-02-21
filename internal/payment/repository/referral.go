package repository

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReferralRepository interface {
	CreateReferralCode(userID string) (string, error)
	ApplyReferralCode(userID, referralCode string) error
	GetReferralBonus(userID string) (float64, error)
}

type referralRepository struct {
	db *sqlx.DB
}

func NewReferralRepository(db *sqlx.DB) ReferralRepository {
	rand.Seed(time.Now().UnixNano())
	return &referralRepository{db: db}
}

// CreateReferralCode generates a referral code and stores it in the referrals table.
func (r *referralRepository) CreateReferralCode(userID string) (string, error) {
	code := "REF" + userID + strconv.Itoa(rand.Intn(10000))
	now := time.Now()
	query := `
		INSERT INTO referrals (user_id, referral_code, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (referral_code) DO NOTHING
	`
	_, err := r.db.Exec(query, userID, code, now)
	if err != nil {
		return "", fmt.Errorf("failed to create referral code: %v", err)
	}
	return code, nil
}

// ApplyReferralCode checks if the referral code exists and, if so, records its usage.
func (r *referralRepository) ApplyReferralCode(userID, referralCode string) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM referrals WHERE referral_code = $1)`
	if err := r.db.Get(&exists, checkQuery, referralCode); err != nil {
		return fmt.Errorf("failed to check referral code: %v", err)
	}
	if !exists {
		return fmt.Errorf("invalid referral code: %s", referralCode)
	}

	// Record the usage and award a bonus (for example, 5.0) -> should come come from product config
	bonus := 5.0
	insertQuery := `
		INSERT INTO referral_usages (applied_user_id, referral_code, bonus, applied_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(insertQuery, userID, referralCode, bonus, time.Now())
	if err != nil {
		return fmt.Errorf("failed to apply referral code: %v", err)
	}
	return nil
}

// GetReferralBonus sums up all bonus amounts awarded to the given user.
func (r *referralRepository) GetReferralBonus(userID string) (float64, error) {
	var totalBonus float64
	query := `SELECT COALESCE(SUM(bonus), 0) FROM referral_usages WHERE applied_user_id = $1`
	if err := r.db.Get(&totalBonus, query, userID); err != nil {
		return 0, fmt.Errorf("failed to get referral bonus: %v", err)
	}
	return totalBonus, nil
}
