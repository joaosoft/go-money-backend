package gomoney

import (
	"database/sql"

	"github.com/joaosoft/go-manager/service"
	"github.com/satori/go.uuid"
)

// Storage ...
type Storage struct {
	conn gomanager.IDB
}

// NewStorage ...
func NewStorage(connection gomanager.IDB) *Storage {
	return &Storage{
		conn: connection,
	}
}

// GetUser ...
func (storage *Storage) GetUser(userID uuid.UUID) (*User, error) {
	row := storage.conn.Get().QueryRow(`
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
	if result, err := storage.conn.Get().Exec(`
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
	if result, err := storage.conn.Get().Exec(`
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
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM money.users
		WHERE user_id = $1
	`, userID.String()); err != nil {
		return err
	}

	return nil
}

// GetTransaction ...
func (storage *Storage) GetTransaction(transactionID uuid.UUID) (*Transaction, error) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
		    user_id,
			category_id,
			price,
			description,
			date,
			updated_at,
			created_at
		FROM money.users
		WHERE transaction_id = $1
	`, transactionID.String())

	transaction := &Transaction{TransactionID: transactionID}
	if err := row.Scan(
		&transaction.UserID,
		&transaction.CategoryID,
		&transaction.Price,
		&transaction.Description,
		&transaction.Date,
		&transaction.UpdatedAt,
		&transaction.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}

	return transaction, nil
}

// CreateTransaction ...
func (storage *Storage) CreateTransaction(newTransaction *Transaction) (*Transaction, error) {
	if result, err := storage.conn.Get().Exec(`
		INSERT INTO money.transactions(transaction_id, user_id, category_id, price, description, date)
		VALUES($1, $2, $3, $4, $5)
	`, newTransaction.TransactionID.String(), newTransaction.UserID.String(), newTransaction.CategoryID, newTransaction.Price, newTransaction.Description, newTransaction.Date); err != nil {
		return nil, err
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.GetTransaction(newTransaction.TransactionID)
	}

	return nil, nil
}

// UpdateTransaction ...
func (storage *Storage) UpdateTransaction(transaction *Transaction) (*Transaction, error) {
	if result, err := storage.conn.Get().Exec(`
		UPDATE money.transactions SET 
			user_id = $1, 
			category_id = $2, 
			price = $3,
			description = $4,
		  	date = $5
		WHERE transaction_id = $6
	`, transaction.UserID, transaction.CategoryID, transaction.Price, transaction.Description, transaction.Date); err != nil {
		return nil, err
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.GetTransaction(transaction.TransactionID)
	}

	return nil, nil
}

// DeleteTransaction ...
func (storage *Storage) DeleteTransaction(transactionID uuid.UUID) error {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM money.transactions
		WHERE transaction_id = $1
	`, transactionID.String()); err != nil {
		return err
	}

	return nil
}
