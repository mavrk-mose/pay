package services

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/mavrk-mose/pay/pkg/utils/http_client"
	"errors"
	"net/http"
	"time"
	
	"github.com/google/uuid"
	. "github.com/mavrk-mose/pay/internal/model"
	"github.com/mavrk-mose/pay/internal/repository"
)

// these are the methods that the payments will trigger
type PaymentActions interface {
	RecordPaymentEvent(event PaymentEvent) error
	ForwardPaymentOrder(order PaymentOrder) error
	NotifyLedger(order PaymentOrder) error
	NotifyWallet(order PaymentOrder) error

	//scheduled payments
	SchedulePayment(order PaymentOrder, scheduleTime time.Time) error // Schedules a payment for future processing
	GetScheduledPayments(userID string) ([]PaymentOrder, error)       // Retrieves scheduled payments
	CancelScheduledPayment(paymentID string) error                    // Cancels a scheduled payment

	//error handling & logging
	LogFailedPayment(event PaymentEvent, reason string) error // Logs a failed payment event
	RetryFailedPayment(paymentID string) error                // Retries processing a failed payment

	//analytics
	// GetPaymentSummary(userID string) (PaymentSummary, error) // Provides a summary of user's payments
	GetTransactionVolumeByPeriod(period string) (int, error) // Retrieves total transaction volume by period

	// fraud detetion
	FlagSuspiciousPayment(paymentID string) error      // Flags a payment as suspicious
	UnflagPayment(paymentID string) error             // Removes suspicious flag from a payment
	QueryFlaggedPayments(userID string) ([]PaymentOrder, error) // Lists flagged payments

	//cancel & refunds
	RequestRefund(paymentID string) error      // Initiates a refund for a payment
	CancelPayment(paymentID string) error      // Cancels a pending payment
	QueryRefundStatus(paymentID string) error // Queries the status of a refund
}

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
