package ports

import (
	"github.com/google/uuid"
	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
)

type Repository[T any] interface {
	Create(*T) error
	Update(*T) error
	Delete(ID uuid.UUID) error
	FindAll() ([]T, error)
	FindOneByConditions(u *domain.User) (T, error)
}
