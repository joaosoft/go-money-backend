package gomoney

import "github.com/satori/go.uuid"

// IStorage ...
type IStorage interface {
	GetUser(id uuid.UUID) (*User, error)
	CreateUser(newUser *User) (*User, error)
	UpdateUser(newUser *User) (*User, error)
	DeleteUser(id uuid.UUID) error
}

// Interactor ...
type Interactor struct {
	storage IStorage
}

// NewInteractor ...
func NewInteractor(repository IStorage) *Interactor {
	return &Interactor{
		storage: repository,
	}
}

// GetUser ...
func (interactor *Interactor) GetUser(userID uuid.UUID) (*User, error) {
	if user, err := interactor.storage.GetUser(userID); err != nil {
		log.Errorf("error getting user on storage %s", err)
		return nil, err
	} else {
		return user, nil
	}
}

// CreateUser ...
func (interactor *Interactor) CreateUser(newUser *User) (*User, error) {
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
	if user, err := interactor.storage.UpdateUser(newUser); err != nil {
		log.Errorf("error updating user on storage %s", err)
		return nil, err
	} else {
		return user, nil
	}
}

// DeleteUser ...
func (interactor *Interactor) DeleteUser(userID uuid.UUID) error {
	if err := interactor.storage.DeleteUser(userID); err != nil {
		log.Errorf("error deleting user on storage %s", err)
		return err
	}
	return nil
}
