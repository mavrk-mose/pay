package service

import (
	"github.com/your_project/repository"
)

type ReferralService interface {
	GenerateReferralCode(userID string) (string, error)
	ApplyReferralCode(userID, referralCode string) error
	GetReferralBonus(userID string) (float64, error)
}

type referralService struct {
	repo repository.ReferralRepository
}

func NewReferralService(repo repository.ReferralRepository) ReferralService {
	return &referralService{repo: repo}
}

func (r *referralService) GenerateReferralCode(userID string) (string, error) {
	return r.repo.CreateReferralCode(userID)
}

func (r *referralService) ApplyReferralCode(userID, referralCode string) error {
	return r.repo.ApplyReferralCode(userID, referralCode)
}

func (r *referralService) GetReferralBonus(userID string) (float64, error) {
	return r.repo.GetReferralBonus(userID)
}
