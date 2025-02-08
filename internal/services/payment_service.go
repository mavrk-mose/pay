package services

import (
	"encoding/json"
	"errors"
	. "github.com/mavrk-mose/pay/internal/model"
	"github.com/mavrk-mose/pay/internal/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type PaymentActions interface {
	RecordPaymentEvent(event PaymentEvent) error
	ForwardPaymentOrder(order PaymentOrder) error
	NotifyLedger(order PaymentOrder) error
	NotifyWallet(order PaymentOrder) error

	SchedulePayment(order PaymentOrder, scheduleTime time.Time) error // Schedules a payment for future processing
	GetScheduledPayments(userID string) ([]PaymentOrder, error)       // Retrieves scheduled payments
	CancelScheduledPayment(paymentID string) error                    // Cancels a scheduled payment

	LogFailedPayment(event PaymentEvent, reason string) error // Logs a failed payment event
	RetryFailedPayment(paymentID string) error                // Retries processing a failed payment

	GetTransactionVolumeByPeriod(period string) (int, error) // Retrieves total transaction volume by period

	FlagSuspiciousPayment(paymentID string) error               // Flags a payment as suspicious
	UnflagPayment(paymentID string) error                       // Removes suspicious flag from a payment
	QueryFlaggedPayments(userID string) ([]PaymentOrder, error) // Lists flagged payments

	RequestRefund(paymentID string) error     // Initiates a refund for a payment
	CancelPayment(paymentID string) error     // Cancels a pending payment
	QueryRefundStatus(paymentID string) error // Queries the status of a refund
}

type PaymentService struct {
	TransactionRepo    repository.TransactionRepo
	HttpClient         *http.Client
	ExternalPaymentURL string
	GenericHttpClient  *utils.GenericHttpClient
}

func NewPaymentService(transactionRepo repository.TransactionRepo, externalPaymentURL string, logger *zap.Logger) *PaymentService {
	return &PaymentService{
		TransactionRepo:    transactionRepo,
		ExternalPaymentURL: externalPaymentURL,
		GenericHttpClient:  utils.NewGenericHttpClient(logger),
	}
}

type ExternalPaymentResponse struct {
	Status      string `json:"status"`
	ExternalRef string `json:"external_ref"`
}

func (s *PaymentService) InitiatePayment(paymentRequest PaymentIntent) (string, error) {
	extResp, err := s.GenericHttpClient.Post(s.ExternalPaymentURL, paymentRequest, map[string]string{})
	if err != nil {
		return "", err
	}

	var response ExternalPaymentResponse
	if err := json.Unmarshal(*extResp, &response); err != nil {
		return "", err
	}

	if response.Status != "ok" && response.Status != "confirmed" {
		return "", errors.New("external payment initiation failed")
	}

	return response.ExternalRef, nil
}
