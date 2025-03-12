package service

import (
	"github.com/mavrk-mose/pay/internal/payment/repository"
)

type DiscountService interface {
	GetMerchantDiscount(merchantID string) (float64, error)
}

type discountService struct {
	repo repository.DiscountRepository
}

func NewDiscountService(repo repository.DiscountRepository) DiscountService {
	return &discountService{repo: repo}
}

func (s *discountService) GetMerchantDiscount(merchantID string) (float64, error) {
	return s.repo.GetMerchantDiscount(merchantID)
}
