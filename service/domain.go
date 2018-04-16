package gomoney

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// User ...
type User struct {
	UserID      uuid.UUID
	Name        string
	Email       string
	Password    string
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// Transaction ...
type Transaction struct {
	TransactionID uuid.UUID
	UserID        uuid.UUID
	CategoryID    uuid.UUID
	Price         decimal.Decimal
	Description   string
	Date          time.Time
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

// Category ...
type Category struct {
	UserID      uuid.UUID
	CategoryID  uuid.UUID
	Name        string
	Description string
}
