package ports

import "gitlab.com/EDteam/workshop-ai-2024/admin/domain"

type UserUseCase interface {
	GenericUseCase[domain.User]
	Login(email, password string) (domain.User, error)
}
