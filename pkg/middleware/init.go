package middleware

type ApiMiddleware struct {
	apiService ports.ApiService
}

func NewApiMiddleware(s ports.ApiService) *ApiMiddleware {
	//TODO: load the public key here
	return &ApiMiddleware{
		apiService: s,
	}
}
