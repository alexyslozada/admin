package client

import (
	"github.com/google/uuid"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	EDtimer "gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports/out"
)

type UseCase struct {
	repo  out.Repository[domain.Client]
	timer EDtimer.Timer
}

func NewUseCase(repo out.Repository[domain.Client], timer EDtimer.Timer) UseCase {
	return UseCase{repo: repo, timer: timer}
}

func (uc UseCase) Create(c *domain.Client) error {
	ID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	c.ID = ID
	c.CreatedAt = uc.timer.Now().Unix()

	return uc.repo.Create(c)
}

func (uc UseCase) Update(c *domain.Client) error {
	c.UpdatedAt = uc.timer.Now().Unix()

	return uc.repo.Update(c)
}

func (uc UseCase) Delete(ID uuid.UUID) error {
	return uc.repo.Delete(ID)
}

func (uc UseCase) FindAll(filters []urler.Filter) ([]domain.Client, error) {
	return uc.repo.FindAll(filters)
}

func (uc UseCase) FindOneByConditions(c *domain.Client) (domain.Client, error) {
	return uc.repo.FindOneByConditions(c)
}
