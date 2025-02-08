package services

import (
	. "github.com/mavrk-mose/pay/internal/ports"
)

type ApiHandler struct {
	store ApiStore
}

func NewApiHandler(store ApiStore) *ApiHandler {
	return &ApiHandler{
		store: store,
	}
}
