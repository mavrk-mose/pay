package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/your_project/models"
)

type VoucherRepository interface {
	CreateVoucher(userID string, amount float64, currency string) (*models.Voucher, error)
	RedeemVoucher(userID, code string) error
}

type voucherRepository struct {
	db *sqlx.DB
}

func NewVoucherRepository(db *sqlx.DB) VoucherRepository {
	return &voucherRepository{db: db}
}

// CreateVoucher generates a voucher code, stores it in the database, and returns the created voucher.
func (r *voucherRepository) CreateVoucher(userID string, amount float64, currency string) (*models.Voucher, error) {
	// Generate a simple voucher code. In production, you might use a more robust method.
	code := fmt.Sprintf("VOUCHER%s%d", userID, time.Now().UnixNano()%10000)
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // Voucher valid for 7 days

	voucher := &models.Voucher{}
	query := `
		INSERT INTO vouchers (user_id, code, amount, currency, expires_at, redeemed, created_at)
		VALUES ($1, $2, $3, $4, $5, false, NOW())
		RETURNING id, user_id, code, amount, currency, expires_at, redeemed, created_at
	`
	err := r.db.Get(voucher, query, userID, code, amount, currency, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create voucher: %v", err)
	}
	return voucher, nil
}

// RedeemVoucher marks the voucher as redeemed if it exists, is not already redeemed, and has not expired.
func (r *voucherRepository) RedeemVoucher(userID, code string) error {
	var voucher models.Voucher
	query := `
		SELECT id, user_id, code, amount, currency, expires_at, redeemed, created_at
		FROM vouchers
		WHERE code = $1 AND user_id = $2
	`
	err := r.db.Get(&voucher, query, code, userID)
	if err != nil {
		return fmt.Errorf("voucher not found: %v", err)
	}
	if voucher.Redeemed {
		return fmt.Errorf("voucher already redeemed")
	}
	if time.Now().After(voucher.ExpiresAt) {
		return fmt.Errorf("voucher expired")
	}

	updateQuery := `UPDATE vouchers SET redeemed = true WHERE id = $1`
	_, err = r.db.Exec(updateQuery, voucher.ID)
	if err != nil {
		return fmt.Errorf("failed to redeem voucher: %v", err)
	}
	return nil
}
