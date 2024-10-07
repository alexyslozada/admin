package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	ProductSubscription = "subscription"
	ProductCourse       = "course"
)

type Sale struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Product        string    `json:"product"`
	ClientID       uuid.UUID `json:"client_id"`
	Client         Client    `json:"client" gorm:"foreignKey:ClientID;references:ID"`
	DateInvoice    int64     `json:"date_invoice"`
	Amount         float64   `json:"amount"`
	IsSubscription bool      `json:"is_subscription"`
	Months         uint8     `json:"months"`
	CreatedAt      int64     `json:"created_at"`
	UpdatedAt      int64     `json:"updated_at"`
	DeletedAt      int64     `json:"deleted_at"`
}

func (Sale) TableName() string {
	return "sales"
}

type SaleSummarized struct {
	Product  string  `json:"product"`
	Amount   float64 `json:"amount"`
	Year     int     `json:"year"`
	Month    int     `json:"month"`
	Quantity int     `json:"quantity"`
}

type SaleSummarizedByClient struct {
	Name     string  `json:"name"`
	Product  string  `json:"product"`
	Amount   float64 `json:"amount"`
	Year     int     `json:"year"`
	Month    int     `json:"month"`
	Quantity int     `json:"quantity"`
}

type SaleResponse struct {
	ID             uuid.UUID `json:"id"`
	Product        string    `json:"product"`
	ClientID       uuid.UUID `json:"client_id"`
	Client         Client    `json:"client"`
	DateInvoice    time.Time `json:"date_invoice"`
	Amount         float64   `json:"amount"`
	IsSubscription bool      `json:"is_subscription"`
	Months         uint8     `json:"months"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}
