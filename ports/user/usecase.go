package user

import "gitlab.com/EDteam/workshop-ai-2024/admin/domain"

type PortUseCase interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
	Login(email, password string) (domain.User, error)
}
