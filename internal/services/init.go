package services

import (
	. "github.com/mavrk-mose/pay/api/internal/ports"
)

type ApiHandler struct {
	store ApiStore
}

func NewApiHandler(store ApiStore) *ApiHandler {
	return &ApiHandler{
		store: store,
	}
}