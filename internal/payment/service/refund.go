package service

import (
	"errors"
	. "github.com/mavrk-mose/pay/internal/payment/models"
	. "github.com/mavrk-mose/pay/internal/payment/repository"
)

type RefundService struct {
	repo *PaymentRepo
}

func NewRefundService(repo *PaymentRepo) *RefundService {
	return &RefundService{repo: repo}
}

func (s *RefundService) ProcessRefund(refund *Refund) error {
	// Business logic for refunds
	if refund.Amount <= 0 {
		return errors.New("invalid refund amount")
	}

	// Save refund as INITIATED
	refund.Status = "INITIATED"
	err := s.repo.CreateRefund(refund)
	if err != nil {
		return err
	}

	// Fake: Assume refund success
	s.repo.UpdateRefundStatus(refund.ID, "COMPLETED")

	return nil
}
