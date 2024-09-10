package ports

import (
	"github.com/google/uuid"
)

type GenericUseCase[T any] interface {
	Create(c *T) error
	Update(c *T) error
	Delete(ID uuid.UUID) error
	FindAll() ([]T, error)
	FindOneByConditions(c *T) (T, error)
}
