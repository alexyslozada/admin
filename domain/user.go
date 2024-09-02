package domain

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `gorm:"type:varchar(100);unique_index"`
	Password  string    `gorm:"type:varchar(100)"`
	Name      string    `gorm:"type:varchar(100);unique_index"`
	CreatedAt int64     `gorm:"type:bigint"`
	UpdatedAt int64     `gorm:"type:bigint"`
	DeletedAt int64     `gorm:"type:bigint"`
}

// TableName Set name of table for gorm
func (User) TableName() string {
	return "users"
}
