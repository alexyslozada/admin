package di

import (
	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/inbound/http"
	"gitlab.com/EDteam/workshop-ai-2024/admin/application/login"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
)

func InitLogin() http.LoginHandler {
	return http.NewLoginHandler(InitLoginUseCase())
}

func InitLoginUseCase() login.UseCase {
	return login.NewUseCase(InitUserUseCase(), timer.NewRealTimer())
}
