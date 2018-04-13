package gomoney

import (
	"time"

	"github.com/satori/go.uuid"
)

type User struct {
	UserID      uuid.UUID
	Name        string
	Email       string
	Password    string
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}
