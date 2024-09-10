package ports

import (
	"github.com/google/uuid"
)

type Repository[T any] interface {
	Create(*T) error
	Update(*T) error
	Delete(ID uuid.UUID) error
	FindAll(preloads ...string) ([]T, error)
	FindOneByConditions(t *T) (T, error)
}
