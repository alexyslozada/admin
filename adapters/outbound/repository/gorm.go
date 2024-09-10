package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gorm[T any] struct {
	db *gorm.DB
}

func NewGorm[T any](db *gorm.DB) Gorm[T] {
	return Gorm[T]{db: db}
}

func (g Gorm[T]) Create(model *T) error {
	return g.db.Create(model).Error
}

func (g Gorm[T]) Update(model *T) error {
	return g.db.Save(model).Error
}

func (g Gorm[T]) Delete(ID uuid.UUID) error {
	var model T
	return g.db.Delete(&model, ID).Error
}

func (g Gorm[T]) FindAll(preloads ...string) ([]T, error) {
	var models []T
	// Preload all relationships
	for _, preload := range preloads {
		g.db = g.db.Preload(preload)
	}
	err := g.db.Find(&models).Error
	return models, err
}

func (g Gorm[T]) FindOneByConditions(model *T) (T, error) {
	err := g.db.Where(model).First(model).Error
	return *model, err
}
