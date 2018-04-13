package gomoney

import (
	"database/sql"

	"github.com/satori/go.uuid"
)

// Storage ...
type Storage struct {
	connection *sql.DB
}

// NewStorage ...
func NewStorage(connection *sql.DB) *Storage {
	return &Storage{
		connection: connection,
	}
}

// GetUser ...
func (storage *Storage) GetUser(userID uuid.UUID) (*User, error) {
	row := storage.connection.QueryRow(`
	    SELECT
		    name,
			email,
			password,
			description,
			updated_at,
			created_at
		FROM money.users
		WHERE user_id = $1
	`, userID.String())

	user := &User{UserID: userID}
	if err := row.Scan(
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Description,
		&user.UpdatedAt,
		&user.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}

	return user, nil
}

// CreateUser ...
func (storage *Storage) CreateUser(newUser *User) (*User, error) {
	if result, err := storage.connection.Exec(`
		INSERT INTO money.users(user_id, name, email, password, description)
		VALUES($1, $2, $3, $4, $5)
	`, newUser.UserID.String(), newUser.Name, newUser.Email, newUser.Password, newUser.Description); err != nil {
		return nil, err
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.GetUser(newUser.UserID)
	}

	return nil, nil
}

// UpdateUser ...
func (storage *Storage) UpdateUser(user *User) (*User, error) {
	if result, err := storage.connection.Exec(`
		UPDATE money.users SET 
			name = $1, 
			email = $2, 
			password = $3,
			description = $4
		WHERE user_id = $5
	`, user.Name, user.Email, user.Password, user.Description, user.UserID.String()); err != nil {
		return nil, err
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.GetUser(user.UserID)
	}

	return nil, nil
}

// DeleteUser ...
func (storage *Storage) DeleteUser(userID uuid.UUID) error {
	if _, err := storage.connection.Exec(`
	    DELETE 
		FROM money.users
		WHERE user_id = $1
	`, userID.String()); err != nil {
		return err
	}

	return nil
}
