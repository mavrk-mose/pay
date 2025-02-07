package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/mavrk-mose/pay/pkg/utils/httpclient"
	"errors"
	"net/http"
	"time"
	
	"github.com/google/uuid"
	. "github.com/mavrk-mose/pay/internal/model"
	"github.com/mavrk-mose/pay/internal/repository"
)

type PaymentService struct {
	TransactionRepo    repository.TransactionRepo
	HttpClient         *http.Client
	ExternalPaymentURL string
}

type ExternalPaymentResponse struct {
	Status      string `json:"status"`
	ExternalRef string `json:"external_ref"`
}


func (s *PaymentService) InitiatePayment(ctx context.Context, paymentRequest PaymentIntent) (string, error) {
	externalRef := uuid.New().String()

	extResp, err := s.GenericHttpClient.Post[PaymentInitiateRequest, ExternalPaymentResponse](s.ExternalPaymentURL, paymentRequest, map[string]string{})
	if err != nil {
		return "", err
	}

	// Optionally, check the external API response status.
	if extResp.Status != "ok" && extResp.Status != "confirmed" {
		return "", errors.New("external payment initiation failed")
	}

	// Return the external reference for tracking purposes.
	return extResp.ExternalRef, nil
}
