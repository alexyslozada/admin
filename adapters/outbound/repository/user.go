package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	return User{db: db}
}

func (u User) Create(user *domain.User) error {
	return u.db.Create(user).Error
}

func (u User) Update(user *domain.User) error {
	return u.db.Save(user).Error
}

func (u User) Delete(ID uuid.UUID) error {
	return u.db.Delete(&domain.User{}, ID).Error
}

func (u User) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := u.db.Find(&users).Error
	return users, err
}

func (u User) FindOneByConditions(model *domain.User) (domain.User, error) {
	err := u.db.Where(model).First(model).Error
	return *model, err
}
