package di

import (
	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/inbound/http"
	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/outbound/repository"
	"gitlab.com/EDteam/workshop-ai-2024/admin/application/client"
	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/infrastructure/db"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
)

func InitClientUseCase() client.UseCase {
	dbGlobal, err := db.NewGorm()
	if err != nil {
		panic(err)
	}

	repo := repository.NewGorm[domain.Client](dbGlobal)
	return client.NewUseCase(repo, timer.NewRealTimer())
}

func InitClientHandler() http.ClientHandler {
	clientUseCase := InitClientUseCase()
	return http.NewClientHandler(clientUseCase)
}
