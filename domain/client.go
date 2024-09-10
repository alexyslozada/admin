package domain

import "github.com/google/uuid"

type Client struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email" gorm:"type:varchar(100);unique"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
	DeletedAt int64     `json:"deleted_at"`
}

func (Client) TableName() string {
	return "clients"
}
