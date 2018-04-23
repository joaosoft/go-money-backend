package gomoney

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// user ...
type user struct {
	UserID       uuid.UUID
	Name         string
	Email        string
	Password     string
	PasswordHash string
	Description  string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

// wallet ...
type wallet struct {
	WalletID    uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Password    string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// image ...
type image struct {
	ImageID     uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Url         string
	FileName    string
	Format      string
	RawImage    []byte
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// transaction ...
type transaction struct {
	TransactionID uuid.UUID
	UserID        uuid.UUID
	WalletID      uuid.UUID
	CategoryID    uuid.UUID
	Price         decimal.Decimal
	Description   string
	Date          time.Time
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

// category ...
type category struct {
	CategoryID  uuid.UUID
	UserID      uuid.UUID
	ImageID     uuid.UUID
	Name        string
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}
