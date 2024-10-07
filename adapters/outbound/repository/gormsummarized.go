package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
)

type GormSummarized[T any] struct {
	db *gorm.DB
}

func NewGormSummarized[T any](db *gorm.DB) GormSummarized[T] {
	return GormSummarized[T]{db: db}
}

func (g GormSummarized[T]) Create(_ *T) error {
	return errors.New("not implemented")
}

func (g GormSummarized[T]) Update(_ *T) error {
	return errors.New("not implemented")
}

func (g GormSummarized[T]) Delete(_ uuid.UUID) error {
	return errors.New("not implemented")
}

func (g GormSummarized[T]) FindAll(filters []urler.Filter, preloads ...string) ([]T, error) {
	var models []T
	// Preload all relationships
	for _, preload := range preloads {
		g.db = g.db.Preload(preload)
	}

	queryRaw := `
		SELECT 
			product, 
			SUM(amount) AS amount, 
			EXTRACT(YEAR FROM to_timestamp(date_invoice)) AS year, 
			EXTRACT(MONTH FROM to_timestamp(date_invoice)) AS month, 
			COUNT(*) AS quantity 
		FROM sales 
		GROUP BY product, year, month
	`
	query := g.db.Raw(queryRaw)
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

	err := query.Scan(&models).Error
	return models, err
}

func (g GormSummarized[T]) FindOneByConditions(model *T) (T, error) {
	// Return the error to implement the method
	return *model, errors.New("not implemented")
}
