package service

import (
	"github.com/your_project/models"
	"github.com/your_project/repository"
)

type VoucherService interface {
	CreateVoucher(userID string, amount float64, currency string) (*models.Voucher, error)
	RedeemVoucher(userID, code string) error
}

type voucherService struct {
	repo repository.VoucherRepository
}

func NewVoucherService(repo repository.VoucherRepository) VoucherService {
	return &voucherService{repo: repo}
}

func (s *voucherService) CreateVoucher(userID string, amount float64, currency string) (*models.Voucher, error) {
	return s.repo.CreateVoucher(userID, amount, currency)
}

func (s *voucherService) RedeemVoucher(userID, code string) error {
	return s.repo.RedeemVoucher(userID, code)
}
