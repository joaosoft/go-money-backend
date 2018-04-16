package gomoney

import "github.com/satori/go.uuid"

// IStorage ...
type IStorage interface {
	GetUser(userID uuid.UUID) (*User, error)
	CreateUser(newUser *User) (*User, error)
	UpdateUser(newUser *User) (*User, error)
	DeleteUser(userID uuid.UUID) error

	GetTransaction(transactionID uuid.UUID) (*Transaction, error)
	CreateTransaction(newTransaction *Transaction) (*Transaction, error)
	UpdateTransaction(newTransaction *Transaction) (*Transaction, error)
	DeleteTransaction(transactionID uuid.UUID) error
}

// Interactor ...
type Interactor struct {
	storage IStorage
	config  *AppConfig
}

// NewInteractor ...
func NewInteractor(repository IStorage, config *AppConfig) *Interactor {
	return &Interactor{
		storage: repository,
		config:  config,
	}
}

// GetUser ...
func (interactor *Interactor) GetUser(userID uuid.UUID) (*User, error) {
	log.WithFields(map[string]interface{}{"method": "GetUser"})
	log.Infof("getting user %s", userID.String())
	if user, err := interactor.storage.GetUser(userID); err != nil {
		log.Errorf("error getting user on storage %s", err)
		return nil, err
	} else {
		return user, nil
	}
}

// CreateUser ...
func (interactor *Interactor) CreateUser(newUser *User) (*User, error) {
	log.WithFields(map[string]interface{}{"method": "CreateUser"})
	log.Infof("creating user %s", newUser.UserID.String())
	newUser.UserID = uuid.NewV4()

	if user, err := interactor.storage.CreateUser(newUser); err != nil {
		log.Errorf("error creating user on storage %s", err)
		return nil, err
	} else {
		return user, nil
	}
}

// UpdateUser ...
func (interactor *Interactor) UpdateUser(newUser *User) (*User, error) {
	log.WithFields(map[string]interface{}{"method": "UpdateUser"})
	log.Infof("updating user %s", newUser.UserID.String())
	if user, err := interactor.storage.UpdateUser(newUser); err != nil {
		log.Errorf("error updating user on storage %s", err)
		return nil, err
	} else {
		return user, nil
	}
}

// DeleteUser ...
func (interactor *Interactor) DeleteUser(userID uuid.UUID) error {
	log.WithFields(map[string]interface{}{"method": "DeleteUser"})
	log.Infof("deleting user %s", userID.String())
	if err := interactor.storage.DeleteUser(userID); err != nil {
		log.Errorf("error deleting user on storage %s", err)
		return err
	}
	return nil
}

// GetTransaction ...
func (interactor *Interactor) GetTransaction(transactionID uuid.UUID) (*Transaction, error) {
	log.WithFields(map[string]interface{}{"method": "GetTransaction"})
	log.Infof("getting transaction %s", transactionID.String())
	if transaction, err := interactor.storage.GetTransaction(transactionID); err != nil {
		log.Errorf("error getting transaction on storage %s", err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// CreateTransaction ...
func (interactor *Interactor) CreateTransaction(newTransaction *Transaction) (*Transaction, error) {
	log.WithFields(map[string]interface{}{"method": "CreateTransaction"})
	log.Infof("creating transaction %s", newTransaction.UserID.String())
	newTransaction.UserID = uuid.NewV4()

	if transaction, err := interactor.storage.CreateTransaction(newTransaction); err != nil {
		log.Errorf("error creating transaction on storage %s", err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// UpdateTransaction ...
func (interactor *Interactor) UpdateTransaction(newTransaction *Transaction) (*Transaction, error) {
	log.WithFields(map[string]interface{}{"method": "UpdateTransaction"})
	log.Infof("updating transaction %s", newTransaction.UserID.String())
	if transaction, err := interactor.storage.UpdateTransaction(newTransaction); err != nil {
		log.Errorf("error updating transaction on storage %s", err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// DeleteTransacgtion ...
func (interactor *Interactor) DeleteTransaction(transactionID uuid.UUID) error {
	log.WithFields(map[string]interface{}{"method": "DeleteTransaction"})
	log.Infof("deleting transaction %s", transactionID.String())
	if err := interactor.storage.DeleteTransaction(transactionID); err != nil {
		log.Errorf("error deleting transaction on storage %s", err)
		return err
	}
	return nil
}
