package di

import "gitlab.com/EDteam/workshop-ai-2024/admin/adapters/inbound/http"

func InitMiddleware() http.Middleware {
	return http.NewMiddleware(InitLoginUseCase())
}
