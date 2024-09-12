package http

import (
	"net/http"

	"gitlab.com/EDteam/workshop-ai-2024/admin/ports"
)

type Middleware struct {
	loginUseCase ports.LoginUseCase
}

func NewMiddleware(loginUseCase ports.LoginUseCase) Middleware {
	return Middleware{loginUseCase: loginUseCase}
}

func (m Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

		isValid, err := m.loginUseCase.ValidateToken(token)
		if !isValid || err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
