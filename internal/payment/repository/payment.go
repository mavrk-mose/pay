package repository

import (
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/payment/models"
)

type PaymentRepo struct {
	DB *sqlx.DB
}

func (r PaymentRepo) CreateRefund(refund *Refund) error {
	return nil
}

func (r PaymentRepo) UpdateRefundStatus(id any, s string) {

}