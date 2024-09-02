package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	EDtimer "gitlab.com/EDteam/workshop-ai-2024/admin/internal/timer"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports"
)

type UseCase struct {
	repo  ports.Repository[domain.User]
	timer EDtimer.Timer
}

func New(repo ports.Repository[domain.User], timer EDtimer.Timer) UseCase {
	return UseCase{repo: repo, timer: timer}
}

func (uc UseCase) Create(user *domain.User) error {
	ID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	user.ID = ID
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)

	user.CreatedAt = uc.timer.Now().Unix()

	return uc.repo.Create(user)
}

func (uc UseCase) Update(user *domain.User) error {
	user.UpdatedAt = uc.timer.Now().Unix()

	return uc.repo.Update(user)
}

func (uc UseCase) Login(email, password string) (domain.User, error) {
	u := domain.User{Email: email}
	u, err := uc.repo.FindOneByConditions(&u)
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, errors.New("invalid password")
	}

	u.Password = ""

	return u, nil
}
