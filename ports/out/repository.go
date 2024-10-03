package out

import (
	"github.com/google/uuid"

	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
)

type Repository[T any] interface {
	Create(*T) error
	Update(*T) error
	Delete(ID uuid.UUID) error
	FindAll(filters []urler.Filter, preloads ...string) ([]T, error)
	FindOneByConditions(t *T) (T, error)
}
