package middleware

import "github.com/mavrk-mose/pay/api/internal/ports"

type ApiMiddleware struct {
	apiService ports.ApiService
}

func NewApiMiddleware(s ports.ApiService) *ApiMiddleware {
	return &ApiMiddleware{
		apiService: s,
	}
}