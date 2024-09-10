package di

import (
	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/inbound/http"
	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/outbound/repository"
	"gitlab.com/EDteam/workshop-ai-2024/admin/application/sale"
	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/infrastructure/db"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
)

func InitSaleUseCase() sale.UseCase {
	dbGlobal, err := db.NewGorm()
	if err != nil {
		panic(err)
	}

	repo := repository.NewGorm[domain.Sale](dbGlobal)
	return sale.NewUseCase(repo, timer.NewRealTimer())
}

func InitSaleHandler() http.SaleHandler {
	saleUseCase := InitSaleUseCase()
	return http.NewSaleHandler(saleUseCase)
}
