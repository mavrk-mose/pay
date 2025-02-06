package model

type PaymentEvent struct {
  EventID     string
  UserID      string
  TotalAmount float64
  Orders      []PaymentOrder
}