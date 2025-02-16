package handler

import (
	. "github.com/mavrk-mose/pay/internal/payment/service"
)

type PaymentHandler struct {
	service *PaymentService
}

func NewPaymentHandler(service *PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}
