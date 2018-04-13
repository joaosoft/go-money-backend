package gomoney

// IStorage ...
type IStorage interface {
	GetUser(id string) (string, error)
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

// Get ...
func (interactor *Interactor) Get(id string) (string, error) {
	return "", nil
}
