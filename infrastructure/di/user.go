package di

import (
	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/outbound/repository"
	"gitlab.com/EDteam/workshop-ai-2024/admin/application/user"
	"gitlab.com/EDteam/workshop-ai-2024/admin/infrastructure/db"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
)

func InitUserUseCase() user.UseCase {
	dbGlobal, err := db.NewGorm()
	if err != nil {
		panic(err)
	}

	repo := repository.NewUser(dbGlobal)
	return user.New(repo, timer.NewRealTimer())
}
