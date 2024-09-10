package sale

import (
	"github.com/google/uuid"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	EDtimer "gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports"
)

type UseCase struct {
	repo  ports.Repository[domain.Sale]
	timer EDtimer.Timer
}

func NewUseCase(repo ports.Repository[domain.Sale], timer EDtimer.Timer) UseCase {
	return UseCase{repo: repo, timer: timer}
}

func (uc UseCase) Create(s *domain.Sale) error {
	ID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	s.ID = ID
	s.CreatedAt = uc.timer.Now().Unix()

	return uc.repo.Create(s)
}

func (uc UseCase) Update(s *domain.Sale) error {
	s.UpdatedAt = uc.timer.Now().Unix()

	return uc.repo.Update(s)
}

func (uc UseCase) Delete(ID uuid.UUID) error {
	return uc.repo.Delete(ID)
}

func (uc UseCase) FindAll() ([]domain.Sale, error) {
	return uc.repo.FindAll("Client")
}

func (uc UseCase) FindOneByConditions(s *domain.Sale) (domain.Sale, error) {
	return uc.repo.FindOneByConditions(s)
}
