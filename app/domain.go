package gomoney

import (
	"time"

	"github.com/shopspring/decimal"
)

// user ...
type user struct {
	UserID      string
	Name        string
	Email       string
	Password    string
	Token       string
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// session ...
type session struct {
	SessionID   string
	UserID      string
	Original    string
	Token       string
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// wallet ...
type wallet struct {
	WalletID    string
	UserID      string
	Name        string
	Description string
	Password    string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// image ...
type image struct {
	ImageID     string
	UserID      string
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
	TransactionID string
	UserID        string
	WalletID      string
	CategoryID    string
	Price         decimal.Decimal
	Description   string
	Date          time.Time
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

// category ...
type category struct {
	CategoryID  string
	UserID      string
	ImageID     string
	Name        string
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}
