package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
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

func (g Gorm[T]) FindAll(filters []urler.Filter, preloads ...string) ([]T, error) {
	var model T
	var models []T
	// Preload all relationships
	for _, preload := range preloads {
		g.db = g.db.Preload(preload)
	}

	query := g.db.Model(&model)
	for _, filter := range filters {
		switch filter.Operator {
		case urler.Equal:
			query = query.Where(filter.Field+" = ?", filter.Value)
		case urler.NotEqual:
			query = query.Where(filter.Field+" != ?", filter.Value)
		case urler.GreaterThan:
			query = query.Where(filter.Field+" > ?", filter.Value)
		case urler.GreaterThanOrEqual:
			query = query.Where(filter.Field+" >= ?", filter.Value)
		case urler.LessThan:
			query = query.Where(filter.Field+" < ?", filter.Value)
		case urler.LessThanOrEqual:
			query = query.Where(filter.Field+" <= ?", filter.Value)
		case urler.Like:
			query = query.Where(filter.Field+" LIKE ?", "%"+filter.Value+"%")
		case urler.In:
			query = query.Where(filter.Field+" IN (?)", filter.Value)
		case urler.NotIn:
			query = query.Where(filter.Field+" NOT IN (?)", filter.Value)
		default:
			query = query.Where(filter.Field+" = ?", filter.Value)
		}
	}

	err := query.Find(&models).Error
	return models, err
}

func (g Gorm[T]) FindOneByConditions(model *T) (T, error) {
	err := g.db.Where(model).First(model).Error
	return *model, err
}
