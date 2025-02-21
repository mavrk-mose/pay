package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type DiscountRepository interface {
	GetMerchantDiscount(merchantID string) (float64, error)
}

type discountRepository struct {
	db *sqlx.DB
}

func NewDiscountRepository(db *sqlx.DB) DiscountRepository {
	return &discountRepository{db: db}
}

// GetMerchantDiscount returns the discount percentage if a valid discount exists.
func (r *discountRepository) GetDiscount(userId string, discountType string) (float64, error) {
	var discount float64
	query := `
		SELECT discount_percentage
		FROM discounts
		WHERE user_id = $1  AND discount_type = $2 AND valid_from <= $3 AND valid_until >= $3
		LIMIT 1
	`
	now := time.Now()
	err := r.db.Get(&discount, query, userId, discountType, now)
	if err != nil {
		return 0, fmt.Errorf("failed to get discount for merchant %s: %v", merchantID, err)
	}
	return discount, nil
}
