package service

import (
	"errors"
	"payment-system/internal/payment/models"
	"payment-system/internal/payment/repository"
)

type RefundService struct {
	repo *repository.PaymentRepository
}

func NewRefundService(repo *repository.PaymentRepository) *RefundService {
	return &RefundService{repo: repo}
}

func (s *RefundService) ProcessRefund(refund *models.Refund) error {
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
