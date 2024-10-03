package http

import (
	"net/http"
	"strings"

	"gitlab.com/EDteam/workshop-ai-2024/admin/ports/app"
)

type Middleware struct {
	loginUseCase app.LoginUseCase
}

func NewMiddleware(loginUseCase app.LoginUseCase) Middleware {
	return Middleware{loginUseCase: loginUseCase}
}

func (m Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
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
