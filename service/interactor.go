package gomoney

import "github.com/satori/go.uuid"

// iStorage ...
type iStorage interface {
	getUsers() ([]*user, error)
	getUser(userID uuid.UUID) (*user, error)
	createUser(newUser *user) (*user, error)
	updateUser(updUser *user) (*user, error)
	deleteUser(userID uuid.UUID) error

	getWallets(userID uuid.UUID) ([]*wallet, error)
	getWallet(userID uuid.UUID, walletID uuid.UUID) (*wallet, error)
	createWallets(newWallets []*wallet) ([]*wallet, error)
	updateWallet(updCategory *wallet) (*wallet, error)
	deleteWallet(userID uuid.UUID, walletID uuid.UUID) error

	getImages(userID uuid.UUID) ([]*image, error)
	getImage(userID uuid.UUID, imageID uuid.UUID) (*image, error)
	createImages(newImage []*image) ([]*image, error)
	updateImage(updImage *image) (*image, error)
	deleteImage(userID uuid.UUID, imageID uuid.UUID) error

	getCategories(userID uuid.UUID) ([]*category, error)
	getCategory(userID uuid.UUID, categoryID uuid.UUID) (*category, error)
	createCategories(newCategory []*category) ([]*category, error)
	updateCategory(updCategory *category) (*category, error)
	deleteCategory(userID uuid.UUID, categoryID uuid.UUID) error

	getTransactions(userID uuid.UUID) ([]*transaction, error)
	getTransaction(userID uuid.UUID, walletID uuid.UUID, transactionID uuid.UUID) (*transaction, error)
	createTransactions(newTransaction []*transaction) ([]*transaction, error)
	updateTransaction(updTransaction *transaction) (*transaction, error)
	deleteTransaction(userID uuid.UUID, walletID uuid.UUID, transactionID uuid.UUID) error
}

// interactor ...
type interactor struct {
	storage iStorage
	config  *appConfig
}

// newInteractor ...
func newInteractor(repository iStorage, config *appConfig) *interactor {
	return &interactor{
		storage: repository,
		config:  config,
	}
}

// getUsers ...
func (interactor *interactor) getUsers() ([]*user, error) {
	log.WithFields(map[string]interface{}{"method": "getUsers"})
	log.Info("getting users")
	if users, err := interactor.storage.getUsers(); err != nil {
		log.Errorf("error getting users on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return users, nil
	}
}

// getUser ...
func (interactor *interactor) getUser(userID uuid.UUID) (*user, error) {
	log.WithFields(map[string]interface{}{"method": "getUser"})
	log.Infof("getting user %s", userID.String())
	if user, err := interactor.storage.getUser(userID); err != nil {
		log.Errorf("error getting user on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return user, nil
	}
}

// createUser ...
func (interactor *interactor) createUser(newUser *user) (*user, error) {
	log.WithFields(map[string]interface{}{"method": "createUser"})

	newUser.UserID = uuid.NewV4()
	log.Infof("creating user %s", newUser.UserID.String())

	if user, err := interactor.storage.createUser(newUser); err != nil {
		log.Errorf("error creating user on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return user, nil
	}
}

// updateUser ...
func (interactor *interactor) updateUser(updUser *user) (*user, error) {
	log.WithFields(map[string]interface{}{"method": "updateUser"})
	log.Infof("updating user %s", updUser.UserID.String())
	if user, err := interactor.storage.updateUser(updUser); err != nil {
		log.Errorf("error updating user on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return user, nil
	}
}

// deleteUser ...
func (interactor *interactor) deleteUser(userID uuid.UUID) error {
	log.WithFields(map[string]interface{}{"method": "deleteUser"})
	log.Infof("deleting user %s", userID.String())
	if err := interactor.storage.deleteUser(userID); err != nil {
		log.Errorf("error deleting user on storage %s", err).ToError(&err)
		return err
	}
	return nil
}

// getWallets ...
func (interactor *interactor) getWallets(userID uuid.UUID) ([]*wallet, error) {
	log.WithFields(map[string]interface{}{"method": "getWallets"})
	log.Infof("getting wallets of user %s", userID.String())
	if wallets, err := interactor.storage.getWallets(userID); err != nil {
		log.Errorf("error getting wallets on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return wallets, nil
	}
}

// getWallet ...
func (interactor *interactor) getWallet(userID uuid.UUID, walletID uuid.UUID) (*wallet, error) {
	log.WithFields(map[string]interface{}{"method": "getWallet"})
	log.Infof("getting wallet %s of user %s", walletID.String(), userID.String())
	if wallet, err := interactor.storage.getWallet(userID, walletID); err != nil {
		log.Errorf("error getting wallet on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return wallet, nil
	}
}

// createWallets ...
func (interactor *interactor) createWallets(newWallets []*wallet) ([]*wallet, error) {
	log.WithFields(map[string]interface{}{"method": "createWallets"})

	log.Info("creating wallets")
	for _, wallet := range newWallets {
		wallet.WalletID = uuid.NewV4()
	}

	if wallets, err := interactor.storage.createWallets(newWallets); err != nil {
		log.Errorf("error creating wallets on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return wallets, nil
	}
}

// updateWallet ...
func (interactor *interactor) updateWallet(updWallet *wallet) (*wallet, error) {
	log.WithFields(map[string]interface{}{"method": "updateWallet"})
	log.Infof("updating wallet %s of user %s", updWallet.UserID.String(), updWallet.UserID.String())
	if wallet, err := interactor.storage.updateWallet(updWallet); err != nil {
		log.Errorf("error updating wallet on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return wallet, nil
	}
}

// deleteWallet ...
func (interactor *interactor) deleteWallet(userID uuid.UUID, walletID uuid.UUID) error {
	log.WithFields(map[string]interface{}{"method": "deleteWallet"})
	log.Infof("deleting wallet %s of user %s", walletID.String(), userID.String())
	if err := interactor.storage.deleteWallet(userID, walletID); err != nil {
		log.Errorf("error deleting wallet on storage %s", err).ToError(&err)
		return err
	}
	return nil
}

// getImages ...
func (interactor *interactor) getImages(userID uuid.UUID) ([]*image, error) {
	log.WithFields(map[string]interface{}{"method": "getImages"})
	log.Infof("getting images of user %s", userID.String())
	if images, err := interactor.storage.getImages(userID); err != nil {
		log.Errorf("error getting images on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return images, nil
	}
}

// getImage ...
func (interactor *interactor) getImage(userID uuid.UUID, imageID uuid.UUID) (*image, error) {
	log.WithFields(map[string]interface{}{"method": "getImage"})
	log.Infof("getting image %s of user %s", imageID.String(), userID.String())
	if image, err := interactor.storage.getImage(userID, imageID); err != nil {
		log.Errorf("error getting image on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return image, nil
	}
}

// createImages ...
func (interactor *interactor) createImages(newImages []*image) ([]*image, error) {
	log.WithFields(map[string]interface{}{"method": "createImages"})

	log.Info("creating images")
	for _, image := range newImages {
		image.ImageID = uuid.NewV4()
	}

	if images, err := interactor.storage.createImages(newImages); err != nil {
		log.Errorf("error creating images on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return images, nil
	}
}

// updateImage ...
func (interactor *interactor) updateImage(updImage *image) (*image, error) {
	log.WithFields(map[string]interface{}{"method": "updateImage"})
	log.Infof("updating image %s of user %s", updImage.UserID.String(), updImage.UserID.String())
	if image, err := interactor.storage.updateImage(updImage); err != nil {
		log.Errorf("error updating image on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return image, nil
	}
}

// deleteImage ...
func (interactor *interactor) deleteImage(userID uuid.UUID, imageID uuid.UUID) error {
	log.WithFields(map[string]interface{}{"method": "deleteImage"})
	log.Infof("deleting image %s of user %s", imageID.String(), userID.String())
	if err := interactor.storage.deleteImage(userID, imageID); err != nil {
		log.Errorf("error deleting image on storage %s", err).ToError(&err)
		return err
	}
	return nil
}

// getCategories ...
func (interactor *interactor) getCategories(userID uuid.UUID) ([]*category, error) {
	log.WithFields(map[string]interface{}{"method": "getCategories"})
	log.Infof("getting categories of user %s", userID.String())
	if categories, err := interactor.storage.getCategories(userID); err != nil {
		log.Errorf("error getting categories on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return categories, nil
	}
}

// getCategory ...
func (interactor *interactor) getCategory(userID uuid.UUID, categoryID uuid.UUID) (*category, error) {
	log.WithFields(map[string]interface{}{"method": "getCategory"})
	log.Infof("getting category %s of user %s", categoryID.String(), userID.String())
	if category, err := interactor.storage.getCategory(userID, categoryID); err != nil {
		log.Errorf("error getting category on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return category, nil
	}
}

// createCategories ...
func (interactor *interactor) createCategories(newCategories []*category) ([]*category, error) {
	log.WithFields(map[string]interface{}{"method": "createCategories"})

	log.Info("creating categories")
	for _, category := range newCategories {
		category.CategoryID = uuid.NewV4()
	}

	if categories, err := interactor.storage.createCategories(newCategories); err != nil {
		log.Errorf("error creating categories on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return categories, nil
	}
}

// updateCategory ...
func (interactor *interactor) updateCategory(updCategory *category) (*category, error) {
	log.WithFields(map[string]interface{}{"method": "updateCategory"})
	log.Infof("updating category %s of user %s", updCategory.UserID.String(), updCategory.UserID.String())
	if category, err := interactor.storage.updateCategory(updCategory); err != nil {
		log.Errorf("error updating category on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return category, nil
	}
}

// deleteCategory ...
func (interactor *interactor) deleteCategory(userID uuid.UUID, categoryID uuid.UUID) error {
	log.WithFields(map[string]interface{}{"method": "deleteCategory"})
	log.Infof("deleting category %s of user %s", categoryID.String(), userID.String())
	if err := interactor.storage.deleteCategory(userID, categoryID); err != nil {
		log.Errorf("error deleting category on storage %s", err).ToError(&err)
		return err
	}
	return nil
}

// getTransaction ...
func (interactor *interactor) getTransactions(userID uuid.UUID) ([]*transaction, error) {
	log.WithFields(map[string]interface{}{"method": "getTransactions"})
	log.Infof("getting transactions of user %s", userID.String())
	if transaction, err := interactor.storage.getTransactions(userID); err != nil {
		log.Errorf("error getting transactions on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// getTransaction ...
func (interactor *interactor) getTransaction(userID uuid.UUID, walletID uuid.UUID, transactionID uuid.UUID) (*transaction, error) {
	log.WithFields(map[string]interface{}{"method": "getTransaction"})
	log.Infof("getting transaction %s of user %s on wallet %s", transactionID.String(), userID.String(), walletID.String())
	if transaction, err := interactor.storage.getTransaction(userID, walletID, transactionID); err != nil {
		log.Errorf("error getting transaction on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// createTransactions ...
func (interactor *interactor) createTransactions(newTransactions []*transaction) ([]*transaction, error) {
	log.WithFields(map[string]interface{}{"method": "createTransactions"})

	log.Info("creating transactions")
	for _, transaction := range newTransactions {
		transaction.TransactionID = uuid.NewV4()
	}

	if transactions, err := interactor.storage.createTransactions(newTransactions); err != nil {
		log.Errorf("error creating transactions on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return transactions, nil
	}
}

// updateTransaction ...
func (interactor *interactor) updateTransaction(updTransaction *transaction) (*transaction, error) {
	log.WithFields(map[string]interface{}{"method": "updateTransaction"})
	log.Infof("updating transaction %s of user %s", updTransaction.UserID.String(), updTransaction.UserID.String())
	if transaction, err := interactor.storage.updateTransaction(updTransaction); err != nil {
		log.Errorf("error updating transaction on storage %s", err).ToError(&err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// deleteTransaction ...
func (interactor *interactor) deleteTransaction(userID uuid.UUID, walletID uuid.UUID, transactionID uuid.UUID) error {
	log.WithFields(map[string]interface{}{"method": "deleteTransaction"})
	log.Infof("deleting transaction %s of user %s", transactionID.String(), userID.String())
	if err := interactor.storage.deleteTransaction(userID, walletID, transactionID); err != nil {
		log.Errorf("error deleting transaction on storage %s", err).ToError(&err)
		return err
	}
	return nil
}
