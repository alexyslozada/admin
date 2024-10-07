package di

import (
	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/inbound/http"
	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/outbound/repository"
	"gitlab.com/EDteam/workshop-ai-2024/admin/application/salesummarized"
	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/infrastructure/db"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
)

func InitSaleSummarizedUseCase() salesummarized.UseCase {
	dbGlobal, err := db.NewGorm()
	if err != nil {
		panic(err)
	}

	repo := repository.NewGormSummarized[domain.SaleSummarized](dbGlobal)
	return salesummarized.NewUseCase(repo, timer.NewRealTimer())
}

func InitSaleSummarizedHandler() http.SaleSummarizedHandler {
	saleSummarizedUseCase := InitSaleSummarizedUseCase()
	return http.NewSaleSummarizedHandler(saleSummarizedUseCase)
}
