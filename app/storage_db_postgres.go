package gomoney

import (
	"database/sql"

	goerror "github.com/joaosoft/go-error/app"
	gomanager "github.com/joaosoft/go-manager/app"
	"github.com/lib/pq"
)

// storagePostgres ...
type storagePostgres struct {
	conn gomanager.IDB
}

// newStoragePostgres ...
func newStoragePostgres(connection gomanager.IDB) *storagePostgres {
	return &storagePostgres{
		conn: connection,
	}
}

// getUsers ...
func (storage *storagePostgres) getUsers() ([]*user, *goerror.ErrorData) {
	rows, err := storage.conn.Get().Query(`
	    SELECT
			user_id,
		    name,
			email,
			password,
			description,
			updated_at,
			created_at
		FROM money.users
	`)
	defer rows.Close()
	if err != nil {
		return nil, goerror.NewError(err)
	}

	users := make([]*user, 0)
	for rows.Next() {
		user := &user{}
		if err := rows.Scan(
			&user.UserID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Token,
			&user.Description,
			&user.UpdatedAt,
			&user.CreatedAt); err != nil {

			if err != sql.ErrNoRows {
				return nil, goerror.NewError(err)
			}
			return nil, nil
		}
		users = append(users, user)
	}

	return users, nil
}

// getUser ...
func (storage *storagePostgres) getUser(userID string) (*user, *goerror.ErrorData) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
		    name,
			email,
			password,
			token,
			description,
			updated_at,
			created_at
		FROM money.users
		WHERE user_id = $1
	`, userID)

	user := &user{UserID: userID}
	if err := row.Scan(
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Token,
		&user.Description,
		&user.UpdatedAt,
		&user.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, goerror.NewError(err)
		}
		return nil, nil
	}

	return user, nil
}

// getUserByEmail ...
func (storage *storagePostgres) getUserByEmail(email string) (*user, *goerror.ErrorData) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
			user_id,
		    name,
			password,
			token,
			description,
			updated_at,
			created_at
		FROM money.users
		WHERE email = $1
	`, email)

	user := &user{Email: email}
	if err := row.Scan(
		&user.UserID,
		&user.Name,
		&user.Password,
		&user.Token,
		&user.Description,
		&user.UpdatedAt,
		&user.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, goerror.NewError(err)
		}
		return nil, nil
	}

	return user, nil
}

// createUser ...
func (storage *storagePostgres) createUser(newUser *user) (*user, *goerror.ErrorData) {
	if result, err := storage.conn.Get().Exec(`
		INSERT INTO money.users(user_id, name, email, password, token, description)
		VALUES($1, $2, $3, $4, $5, $6)
	`, newUser.UserID, newUser.Name, newUser.Email, newUser.Password, newUser.Token, newUser.Description); err != nil {
		return nil, goerror.NewError(err)
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.getUser(newUser.UserID)
	}

	return nil, nil
}

// updateUser ...
func (storage *storagePostgres) updateUser(user *user) (*user, *goerror.ErrorData) {
	if result, err := storage.conn.Get().Exec(`
		UPDATE money.users SET 
			name = $1, 
			email = $2, 
			password = $3,
			token = $4,
			description = $5
		WHERE user_id = $6
	`, user.Name, user.Email, user.Password, user.Token, user.Description, user.UserID); err != nil {
		return nil, goerror.NewError(err)
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.getUser(user.UserID)
	}

	return nil, nil
}

// deleteUser ...
func (storage *storagePostgres) deleteUser(userID string) *goerror.ErrorData {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM money.users
		WHERE user_id = $1
	`, userID); err != nil {
		return goerror.NewError(err)
	}

	return nil
}

// getSessions ...
func (storage *storagePostgres) getSessions(userID string) ([]*session, *goerror.ErrorData) {
	rows, err := storage.conn.Get().Query(`
	    SELECT
			session_id,
		    original,
			token,
			description,
			updated_at,
			created_at
		FROM money.sessions
	  	WHERE user_id = $1
	`, userID)
	defer rows.Close()
	if err != nil {
		return nil, goerror.NewError(err)
	}

	sessions := make([]*session, 0)
	for rows.Next() {
		session := &session{UserID: userID}
		if err := rows.Scan(
			&session.SessionID,
			&session.Original,
			&session.Token,
			&session.Description,
			&session.UpdatedAt,
			&session.CreatedAt); err != nil {

			if err != sql.ErrNoRows {
				return nil, goerror.NewError(err)
			}
			return nil, nil
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// getSession ...
func (storage *storagePostgres) getSession(userID string, token string) (*session, *goerror.ErrorData) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
		    session_id,
			original,
			description,
			updated_at,
			created_at
		FROM money.sessions
		WHERE user_id = $1 AND token = $2
	`, userID, token)

	session := &session{UserID: userID, Token: token}
	if err := row.Scan(
		&session.SessionID,
		&session.Original,
		&session.Description,
		&session.UpdatedAt,
		&session.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, goerror.NewError(err)
		}
		return nil, nil
	}

	return session, nil
}

// createSession ...
func (storage *storagePostgres) createSession(newSession *session) (*session, *goerror.ErrorData) {
	if result, err := storage.conn.Get().Exec(`
		INSERT INTO money.sessions(session_id, user_id, original, token, description)
		VALUES($1, $2, $3, $4, $5)
	`, newSession.SessionID, newSession.UserID, newSession.Original, newSession.Token, newSession.Description); err != nil {
		return nil, goerror.NewError(err)
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.getSession(newSession.UserID, newSession.Token)
	}

	return nil, nil
}

// deleteSession ...
func (storage *storagePostgres) deleteSession(userID string, token string) *goerror.ErrorData {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM money.sessions
		WHERE user_id = $1 AND token = $2
	`, userID, token); err != nil {
		return goerror.NewError(err)
	}

	return nil
}

// deleteSessions ...
func (storage *storagePostgres) deleteSessions(userID string) *goerror.ErrorData {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM money.sessions
		WHERE user_id = $1
	`, userID); err != nil {
		return goerror.NewError(err)
	}

	return nil
}

// getWallets ...
func (storage *storagePostgres) getWallets(userID string) ([]*wallet, *goerror.ErrorData) {
	rows, err := storage.conn.Get().Query(`
	     SELECT
			wallet_id,
			name,
			description,
			password,
			updated_at,
			created_at
		FROM money.wallets
		WHERE user_id = $1
	`, userID)

	defer rows.Close()
	if err != nil {
		return nil, goerror.NewError(err)
	}

	wallets := make([]*wallet, 0)
	for rows.Next() {
		wallet := &wallet{
			UserID: userID,
		}
		if err := rows.Scan(
			&wallet.WalletID,
			&wallet.Name,
			&wallet.Description,
			&wallet.Password,
			&wallet.UpdatedAt,
			&wallet.CreatedAt); err != nil {

			if err != sql.ErrNoRows {
				return nil, goerror.NewError(err)
			}
			return nil, nil
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

// getWallet ...
func (storage *storagePostgres) getWallet(userID string, walletID string) (*wallet, *goerror.ErrorData) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
			name,
			description,
			password,
			updated_at,
			created_at
		FROM money.wallets
		WHERE user_id = $1 AND wallet_id = $2
	`, userID, walletID)

	wallet := &wallet{
		WalletID: walletID,
		UserID:   userID,
	}
	if err := row.Scan(
		&wallet.Name,
		&wallet.Description,
		&wallet.Password,
		&wallet.UpdatedAt,
		&wallet.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, goerror.NewError(err)
		}
		return nil, nil
	}

	return wallet, nil
}

// createWallets ...
func (storage *storagePostgres) createWallets(newWallets []*wallet) ([]*wallet, *goerror.ErrorData) {
	tx, err := storage.conn.Get().Begin()
	if err != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	stmt, errItem := tx.Prepare(pq.CopyInSchema("money", "wallets", "wallet_id", "user_id", "name", "description", "password"))
	if errItem != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	for _, newWallet := range newWallets {
		if _, err := stmt.Exec(newWallet.WalletID, newWallet.UserID, newWallet.Name, newWallet.Description, newWallet.Password); err != nil {
			tx.Rollback()
			return nil, goerror.NewError(err)
		}
	}

	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	if err := stmt.Close(); err != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	tx.Commit()

	// get created wallets
	createdWallets := make([]*wallet, 0)
	for _, newWallet := range newWallets {
		wallet, err := storage.getWallet(newWallet.UserID, newWallet.WalletID)
		if err != nil {
			return nil, goerror.NewError(err)
		}
		createdWallets = append(createdWallets, wallet)
	}

	return createdWallets, nil
}

// updateWallet ...
func (storage *storagePostgres) updateWallet(wallet *wallet) (*wallet, *goerror.ErrorData) {
	if result, err := storage.conn.Get().Exec(`
		UPDATE money.wallets SET 
			name = $1,
			description = $2,
			password = $3
		WHERE user_id = $4 AND wallet_id = $5
	`, wallet.Name, wallet.Description, wallet.Password, wallet.UserID, wallet.WalletID); err != nil {
		return nil, goerror.NewError(err)
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.getWallet(wallet.UserID, wallet.WalletID)
	}

	return nil, nil
}

// deleteWallet ...
func (storage *storagePostgres) deleteWallet(userID string, walletID string) *goerror.ErrorData {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM money.wallets
		WHERE user_id = $1 AND wallet_id = $2
	`, userID, walletID); err != nil {
		return goerror.NewError(err)
	}

	return nil
}

// getImages ...
func (storage *storagePostgres) getImages(userID string) ([]*image, *goerror.ErrorData) {
	rows, err := storage.conn.Get().Query(`
	     SELECT
			image_id,
			name,
			description,
			url,
			file_name,
			format,
			raw_image,
			updated_at,
			created_at
		FROM money.images
		WHERE user_id = $1
	`, userID)

	defer rows.Close()
	if err != nil {
		return nil, goerror.NewError(err)
	}

	images := make([]*image, 0)
	for rows.Next() {
		image := &image{
			UserID: userID,
		}
		if err := rows.Scan(
			&image.ImageID,
			&image.Name,
			&image.Description,
			&image.Url,
			&image.FileName,
			&image.Format,
			&image.RawImage,
			&image.UpdatedAt,
			&image.CreatedAt); err != nil {

			if err != sql.ErrNoRows {
				return nil, goerror.NewError(err)
			}
			return nil, nil
		}
		images = append(images, image)
	}

	return images, nil
}

// getImage ...
func (storage *storagePostgres) getImage(userID string, imageID string) (*image, *goerror.ErrorData) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
			name,
			description,
			url,
			file_name,
			format,
			raw_image,
			updated_at,
			created_at
		FROM money.images
		WHERE user_id = $1 AND image_id = $2
	`, userID, imageID)

	image := &image{
		UserID:  userID,
		ImageID: imageID,
	}
	if err := row.Scan(
		&image.Name,
		&image.Description,
		&image.Url,
		&image.FileName,
		&image.Format,
		&image.RawImage,
		&image.UpdatedAt,
		&image.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, goerror.NewError(err)
		}
		return nil, nil
	}

	return image, nil
}

// createImage ...
func (storage *storagePostgres) createImage(newImage *image) (*image, *goerror.ErrorData) {
	if result, err := storage.conn.Get().Exec(`
		INSERT INTO money.images(image_id, user_id, name, description, url, file_name, format, raw_image)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`, newImage.ImageID, newImage.UserID, newImage.Name, newImage.Description, newImage.Url, newImage.FileName, newImage.Format, newImage.RawImage); err != nil {
		return nil, goerror.NewError(err)
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.getImage(newImage.UserID, newImage.ImageID)
	}
	return nil, nil
}

// updateImage ...
func (storage *storagePostgres) updateImage(updImage *image) (*image, *goerror.ErrorData) {
	if result, err := storage.conn.Get().Exec(`
		UPDATE money.images SET 
			name = $1,
			description = $2,
			url = $3,
			file_name = $4,
			format = $5,
			raw_image = $6
		WHERE user_id = $7 AND image_id = $8
	`, updImage.Name, updImage.Description, updImage.Url, updImage.FileName, updImage.Format, updImage.RawImage, updImage.UserID, updImage.ImageID); err != nil {
		return nil, goerror.NewError(err)
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.getImage(updImage.UserID, updImage.ImageID)
	}

	return nil, nil
}

// deleteImage ...
func (storage *storagePostgres) deleteImage(userID string, imageID string) *goerror.ErrorData {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM money.images
		WHERE user_id = $1 AND image_id = $2
	`, userID, imageID); err != nil {
		return goerror.NewError(err)
	}

	return nil
}

// getCategories ...
func (storage *storagePostgres) getCategories(userID string) ([]*category, *goerror.ErrorData) {
	rows, err := storage.conn.Get().Query(`
	     SELECT
			category_id,
			image_id,
			name,
			description,
			updated_at,
			created_at
		FROM money.categories
		WHERE user_id = $1
	`, userID)

	defer rows.Close()
	if err != nil {
		return nil, goerror.NewError(err)
	}

	categories := make([]*category, 0)
	for rows.Next() {
		category := &category{
			UserID: userID,
		}
		if err := rows.Scan(
			&category.CategoryID,
			&category.ImageID,
			&category.Name,
			&category.Description,
			&category.UpdatedAt,
			&category.CreatedAt); err != nil {

			if err != sql.ErrNoRows {
				return nil, goerror.NewError(err)
			}
			return nil, nil
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// getCategory ...
func (storage *storagePostgres) getCategory(userID string, categoryID string) (*category, *goerror.ErrorData) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
			image_id,
			name,
			description,
			updated_at,
			created_at
		FROM money.categories
		WHERE user_id = $1 AND category_id = $2
	`, userID, categoryID)

	category := &category{
		CategoryID: categoryID,
		UserID:     userID,
	}
	if err := row.Scan(
		&category.ImageID,
		&category.Name,
		&category.Description,
		&category.UpdatedAt,
		&category.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, goerror.NewError(err)
		}
		return nil, nil
	}

	return category, nil
}

// createCategories ...
func (storage *storagePostgres) createCategories(newCategories []*category) ([]*category, *goerror.ErrorData) {
	tx, err := storage.conn.Get().Begin()
	if err != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	stmt, errItem := tx.Prepare(pq.CopyInSchema("money", "categories", "category_id", "user_id", "image_id", "name", "description"))
	if errItem != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	for _, newCategory := range newCategories {
		if _, err := stmt.Exec(newCategory.CategoryID, newCategory.UserID, newCategory.ImageID, newCategory.Name, newCategory.Description); err != nil {
			tx.Rollback()
			return nil, goerror.NewError(err)
		}
	}

	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	if err := stmt.Close(); err != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	tx.Commit()

	// get created categories
	createdCategories := make([]*category, 0)
	for _, newCategory := range newCategories {
		category, err := storage.getCategory(newCategory.UserID, newCategory.CategoryID)
		if err != nil {
			return nil, goerror.NewError(err)
		}
		createdCategories = append(createdCategories, category)
	}

	return createdCategories, nil
}

// updateCategory ...
func (storage *storagePostgres) updateCategory(category *category) (*category, *goerror.ErrorData) {
	if result, err := storage.conn.Get().Exec(`
		UPDATE money.categories SET 
			image_id = $1
			name = $2,
			description = $3,
		WHERE user_id = $4 AND category_id = $5
	`, category.ImageID, category.Name, category.Description, category.UserID, category.CategoryID); err != nil {
		return nil, goerror.NewError(err)
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.getCategory(category.UserID, category.CategoryID)
	}

	return nil, nil
}

// deleteCategory ...
func (storage *storagePostgres) deleteCategory(userID string, categoryID string) *goerror.ErrorData {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM money.categories
		WHERE user_id = $1 AND category_id = $2
	`, userID, categoryID); err != nil {
		return goerror.NewError(err)
	}

	return nil
}

// getTransactions ...
func (storage *storagePostgres) getTransactions(userID string) ([]*transaction, *goerror.ErrorData) {
	rows, err := storage.conn.Get().Query(`
	     SELECT
			wallet_id,
			transaction_id,
			category_id,
			price,
			description,
			date,
			updated_at,
			created_at
		FROM money.transactions
		WHERE user_id = $1
	`, userID)

	defer rows.Close()
	if err != nil {
		return nil, goerror.NewError(err)
	}

	transactions := make([]*transaction, 0)
	for rows.Next() {
		transaction := &transaction{
			UserID: userID,
		}
		if err := rows.Scan(
			&transaction.WalletID,
			&transaction.TransactionID,
			&transaction.CategoryID,
			&transaction.Price,
			&transaction.Description,
			&transaction.Date,
			&transaction.UpdatedAt,
			&transaction.CreatedAt); err != nil {

			if err != sql.ErrNoRows {
				return nil, goerror.NewError(err)
			}
			return nil, nil
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// getTransaction ...
func (storage *storagePostgres) getTransaction(userID string, walletID string, transactionID string) (*transaction, *goerror.ErrorData) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
			category_id,
			price,
			description,
			date,
			updated_at,
			created_at
		FROM money.transactions
		WHERE user_id = $1 AND wallet_id = $2 AND transaction_id = $3
	`, userID, walletID, transactionID)

	transaction := &transaction{
		UserID:        userID,
		WalletID:      walletID,
		TransactionID: transactionID,
	}
	if err := row.Scan(
		&transaction.CategoryID,
		&transaction.Price,
		&transaction.Description,
		&transaction.Date,
		&transaction.UpdatedAt,
		&transaction.CreatedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, goerror.NewError(err)
		}
		return nil, nil
	}

	return transaction, nil
}

// createTransactions ...
func (storage *storagePostgres) createTransactions(newTransactions []*transaction) ([]*transaction, *goerror.ErrorData) {
	tx, err := storage.conn.Get().Begin()
	if err != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	stmt, errItem := tx.Prepare(pq.CopyInSchema("money", "transactions", "transaction_id", "user_id", "wallet_id", "category_id", "price", "description", "date"))
	if errItem != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	for _, newTransaction := range newTransactions {
		if _, err := stmt.Exec(newTransaction.TransactionID, newTransaction.UserID, newTransaction.WalletID, newTransaction.CategoryID, newTransaction.Price, newTransaction.Description, newTransaction.Date); err != nil {
			tx.Rollback()
			return nil, goerror.NewError(err)
		}
	}

	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	if err := stmt.Close(); err != nil {
		tx.Rollback()
		return nil, goerror.NewError(err)
	}

	tx.Commit()

	// get created transactions
	createdTransactions := make([]*transaction, 0)
	for _, newTransaction := range newTransactions {
		transaction, err := storage.getTransaction(newTransaction.UserID, newTransaction.WalletID, newTransaction.TransactionID)
		if err != nil {
			return nil, goerror.NewError(err)
		}
		createdTransactions = append(createdTransactions, transaction)
	}

	return createdTransactions, nil
}

// updateTransaction ...
func (storage *storagePostgres) updateTransaction(transaction *transaction) (*transaction, *goerror.ErrorData) {
	if result, err := storage.conn.Get().Exec(`
		UPDATE money.transactions SET 
			category_id = $1, 
			price = $2,
			description = $3,
		  	date = $4
		WHERE user_id = $5 AND wallet_id = $6 AND transaction_id = $7
	`, transaction.CategoryID, transaction.Price, transaction.Description, transaction.Date, transaction.UserID, transaction.WalletID, transaction.TransactionID); err != nil {
		return nil, goerror.NewError(err)
	} else if rows, _ := result.RowsAffected(); rows > 0 {
		return storage.getTransaction(transaction.UserID, transaction.WalletID, transaction.TransactionID)
	}

	return nil, nil
}

// deleteTransaction ...
func (storage *storagePostgres) deleteTransaction(userID string, walletID string, transactionID string) *goerror.ErrorData {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM money.transactions
		WHERE user_id = $1 AND wallet_id = $2 AND transaction_id = $3
	`, userID, walletID, transactionID); err != nil {
		return goerror.NewError(err)
	}

	return nil
}
