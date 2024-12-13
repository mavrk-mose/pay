package model

type PaymentOrder struct {
  OrderID      string
  Amount       float64
  SellerID     string
  PaymentMethod string
}