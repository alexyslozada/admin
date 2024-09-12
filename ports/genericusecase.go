package ports

import (
	"github.com/google/uuid"

	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
)

type GenericUseCase[T any] interface {
	Create(c *T) error
	Update(c *T) error
	Delete(ID uuid.UUID) error
	FindAll(filters []urler.Filter) ([]T, error)
	FindOneByConditions(c *T) (T, error)
}
