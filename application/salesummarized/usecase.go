package salesummarized

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	EDtimer "gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports/out"
)

type UseCase struct {
	repo  out.Repository[domain.SaleSummarized]
	timer EDtimer.Timer
}

func NewUseCase(repo out.Repository[domain.SaleSummarized], timer EDtimer.Timer) UseCase {
	return UseCase{repo: repo, timer: timer}
}

func (uc UseCase) Create(s *domain.SaleSummarized) error {
	return errors.New("not implemented")
}

func (uc UseCase) Update(s *domain.SaleSummarized) error {
	return errors.New("not implemented")
}

func (uc UseCase) Delete(ID uuid.UUID) error {
	return errors.New("not implemented")
}

func (uc UseCase) FindAll(filters []urler.Filter) ([]domain.SaleSummarized, error) {
	log.Printf("salessummarized.FindAll() filters: %+v", filters)
	// Parse filters to the real business logic
	for i, filter := range filters {
		if strings.EqualFold(filter.Field, "from") {
			filters[i].Field = "date_invoice"
			filters[i].Operator = urler.GreaterThanOrEqual
			from, err := time.Parse(time.DateOnly, filter.Value)
			if err != nil {
				return nil, err
			}
			filters[i].Value = strconv.FormatInt(from.Unix(), 10)
			continue
		}

		if strings.EqualFold(filter.Field, "to") {
			filters[i].Field = "date_invoice"
			filters[i].Operator = urler.LessThanOrEqual
			to, err := time.Parse(time.DateOnly, filter.Value)
			if err != nil {
				return nil, err
			}
			filters[i].Value = strconv.FormatInt(to.Unix(), 10)
			continue
		}

		// Remove the filter if it's not needed
		if len(filters) > 0 {
			filters = append(filters[:i], filters[i+1:]...)
		}
	}

	return uc.repo.FindAll(filters, "Client")
}

func (uc UseCase) FindOneByConditions(s *domain.SaleSummarized) (domain.SaleSummarized, error) {
	return domain.SaleSummarized{}, errors.New("not implemented")
}
