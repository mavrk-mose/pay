package response

import (
	"net/http"
	"time"
)

type ApiResponse[T any] struct {
	TransactionID string `json:"transactionId"`
	Timestamp       string `json:"timestamp"`
	Status         int    `json:"status"`
	Message        string `json:"message"`
	Data          T      `json:"data"`
}

// NewApiResponse creates a new ApiResponse instance.
func NewApiResponse[T any](transactionId string, status int, message string, data T) ApiResponse[T] {
	return ApiResponse[T]{
		TransactionID: transactionId,
		Timestamp:       time.Now().Format(time.RFC3339),
		Status:         status,
		Message:        message,
		Data:           data,
	}
}

func NewErrorResponse[T any](ref string, message string) ApiResponse[T] {
	return ApiResponse[T]{
		ReferenceNumber: ref,
		Timestamp:       time.Now().Format(time.RFC3339),
		Status:          400,
		Message:         message,
		Data:            *new(T), // Zero value for T
	}
}

