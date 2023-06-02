package gomoney

import (
	"fmt"

	"github.com/joaosoft/errors"
)

// iStorageDB ...
type iStorageDB interface {
	getSession(userID string, token string) (*session, error)
	getSessions(userID string) ([]*session, error)
	createSession(newSession *session) (*session, error)
	deleteSession(userID string, token string) error
	deleteSessions(userID string) error

	getUsers() ([]*user, error)
	getUser(userID string) (*user, error)
	getUserByEmail(email string) (*user, error)
	createUser(newUser *user) (*user, error)
	updateUser(updUser *user) (*user, error)
	deleteUser(userID string) error

	getWallets(userID string) ([]*wallet, error)
	getWallet(userID string, walletID string) (*wallet, error)
	createWallets(newWallets []*wallet) ([]*wallet, error)
	updateWallet(updCategory *wallet) (*wallet, error)
	deleteWallet(userID string, walletID string) error

	getImages(userID string) ([]*image, error)
	getImage(userID string, imageID string) (*image, error)
	createImage(newImage *image) (*image, error)
	updateImage(updImage *image) (*image, error)
	deleteImage(userID string, imageID string) error

	getCategories(userID string) ([]*category, error)
	getCategory(userID string, categoryID string) (*category, error)
	createCategories(newCategory []*category) ([]*category, error)
	updateCategory(updCategory *category) (*category, error)
	deleteCategory(userID string, categoryID string) error

	getTransactions(userID string) ([]*transaction, error)
	getTransaction(userID string, walletID string, transactionID string) (*transaction, error)
	createTransactions(newTransaction []*transaction) ([]*transaction, error)
	updateTransaction(updTransaction *transaction) (*transaction, error)
	deleteTransaction(userID string, walletID string, transactionID string) error
}

// iStorageDropbox ...
type iStorageDropbox interface {
	upload(path string, data []byte) error
	download(path string) ([]byte, error)
	delete(path string) error
}

// interactor ...
type interactor struct {
	storageDB      iStorageDB
	storageDropbox iStorageDropbox
	config         *MoneyConfig
}

// newInteractor ...
func newInteractor(storageDB iStorageDB, storageDropbox iStorageDropbox, config *MoneyConfig) *interactor {
	return &interactor{
		storageDB:      storageDB,
		storageDropbox: storageDropbox,
		config:         config,
	}
}

// getUsers ...
func (interactor *interactor) getUsers() ([]*user, error) {
	log.WithFields(map[string]interface{}{"method": "getUsers"})
	log.Info("getting users")
	if users, err := interactor.storageDB.getUsers(); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting users on storage database %s", err)
		return nil, err
	} else {
		return users, nil
	}
}

// getUser ...
func (interactor *interactor) getUser(userID string) (*user, error) {
	log.WithFields(map[string]interface{}{"method": "getUser"})
	log.Infof("getting user %s", userID)
	if user, err := interactor.storageDB.getUser(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting user on storage database %s", err)
		return nil, err
	} else {
		return user, nil
	}
}

// getUserByEmail ...
func (interactor *interactor) getUserByEmail(email string) (*user, error) {
	log.WithFields(map[string]interface{}{"method": "getUserByEmail"})
	log.Infof("getting user by email %s", email)
	if user, err := interactor.storageDB.getUserByEmail(email); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting user by email %s on storage database %s", email, err)
		return nil, err
	} else {
		return user, nil
	}
}

// createUser ...
func (interactor *interactor) createUser(newUser *user) (*user, error) {
	log.WithFields(map[string]interface{}{"method": "createUser"})

	newUser.UserID = genUI()
	passwordToken, err := generateToken(authentication, []byte(newUser.Password))
	if err != nil {
		newErr := errors.New(errors.LevelError, 1, err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr}).
			Error("error when generating password token")
		return nil, newErr
	}
	newUser.Token = passwordToken

	log.Infof("creating user %s", newUser.UserID)

	if user, err := interactor.storageDB.createUser(newUser); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error creating user on storage database %s", err)
		return nil, err
	} else {
		return user, nil
	}
}

// updateUser ...
func (interactor *interactor) updateUser(updUser *user) (*user, error) {
	log.WithFields(map[string]interface{}{"method": "updateUser"})
	log.Infof("updating user %s", updUser.UserID)

	passwordToken, err := generateToken(authentication, []byte(updUser.Password))
	if err != nil {
		newErr := errors.New(errors.LevelError, 1, err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr}).
			Error("error when generating password token")
		return nil, newErr
	}
	updUser.Token = passwordToken

	if user, err := interactor.storageDB.updateUser(updUser); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error updating user on storage database %s", err)
		return nil, err
	} else {
		return user, nil
	}
}

// deleteUser ...
func (interactor *interactor) deleteUser(userID string) error {
	log.WithFields(map[string]interface{}{"method": "deleteUser"})
	log.Infof("deleting user %s", userID)
	if err := interactor.storageDB.deleteUser(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error deleting user on storage database %s", err)
		return err
	}
	return nil
}

// getSessions ...
func (interactor *interactor) getSessions(userID string) ([]*session, error) {
	log.WithFields(map[string]interface{}{"method": "getSessions"})
	log.Info("getting sessions")
	if sessions, err := interactor.storageDB.getSessions(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting sessions on storage database %s", err)
		return nil, err
	} else {
		return sessions, nil
	}
}

// getSession ...
func (interactor *interactor) getSession(userID string, token string) (*session, error) {
	log.WithFields(map[string]interface{}{"method": "getSession"})
	log.Infof("getting session with token %s", userID)
	if session, err := interactor.storageDB.getSession(userID, token); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting session on storage database %s", err)
		return nil, err
	} else {
		return session, nil
	}
}

// createSession ...
func (interactor *interactor) createSession(newSession *session) (*session, error) {
	log.WithFields(map[string]interface{}{"method": "createSession"})

	random := genUI()
	newSession.SessionID = genUI()
	token, err := generateToken(authentication, []byte(random))
	if err != nil {
		newErr := errors.New(errors.LevelError, 1, err)
		log.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error getting sessions on storage database %s", err)
		return nil, newErr
	}
	newSession.Original = random
	newSession.Token = token

	log.Infof("creating session for user %s", newSession.UserID)

	if user, err := interactor.storageDB.createSession(newSession); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error creating session on storage database %s", err)
		return nil, err
	} else {
		return user, nil
	}
}

// deleteSessions ...
func (interactor *interactor) deleteSessions(userID string) error {
	log.WithFields(map[string]interface{}{"method": "deleteSessions"})
	log.Infof("deleting sessions of user %s", userID)
	if err := interactor.storageDB.deleteSessions(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error deleting sessions of user %s on storage database %s", userID, err)
		return err
	}
	return nil
}

// deleteSession ...
func (interactor *interactor) deleteSession(userID string, token string) error {
	log.WithFields(map[string]interface{}{"method": "deleteSession"})
	log.Infof("deleting session with token %s", userID)
	if err := interactor.storageDB.deleteSession(userID, token); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error deleting session of user %s  with token %s on storage database %s", userID, token, err)
		return err
	}
	return nil
}

// getWallets ...
func (interactor *interactor) getWallets(userID string) ([]*wallet, error) {
	log.WithFields(map[string]interface{}{"method": "getWallets"})
	log.Infof("getting wallets of user %s", userID)
	if wallets, err := interactor.storageDB.getWallets(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting wallets on storage database %s", err)
		return nil, err
	} else {
		return wallets, nil
	}
}

// getWallet ...
func (interactor *interactor) getWallet(userID string, walletID string) (*wallet, error) {
	log.WithFields(map[string]interface{}{"method": "getWallet"})
	log.Infof("getting wallet %s of user %s", walletID, userID)
	if wallet, err := interactor.storageDB.getWallet(userID, walletID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting wallet on storage database %s", err)
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
		wallet.WalletID = genUI()
	}

	if wallets, err := interactor.storageDB.createWallets(newWallets); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error creating wallets on storage database %s", err)
		return nil, err
	} else {
		return wallets, nil
	}
}

// updateWallet ...
func (interactor *interactor) updateWallet(updWallet *wallet) (*wallet, error) {
	log.WithFields(map[string]interface{}{"method": "updateWallet"})
	log.Infof("updating wallet %s of user %s", updWallet.UserID, updWallet.UserID)
	if wallet, err := interactor.storageDB.updateWallet(updWallet); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error updating wallet on storage database %s", err)
		return nil, err
	} else {
		return wallet, nil
	}
}

// deleteWallet ...
func (interactor *interactor) deleteWallet(userID string, walletID string) error {
	log.WithFields(map[string]interface{}{"method": "deleteWallet"})
	log.Infof("deleting wallet %s of user %s", walletID, userID)
	if err := interactor.storageDB.deleteWallet(userID, walletID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error deleting wallet on storage database %s", err)
		return err
	}
	return nil
}

// getImages ...
func (interactor *interactor) getImages(userID string) ([]*image, error) {
	log.WithFields(map[string]interface{}{"method": "getImages"})
	log.Infof("getting images of user %s", userID)
	if images, err := interactor.storageDB.getImages(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting images on storage database %s", err)
		return nil, err
	} else {
		return images, nil
	}
}

// getImage ...
func (interactor *interactor) getImage(userID string, imageID string) (*image, error) {
	log.WithFields(map[string]interface{}{"method": "getImage"})
	log.Infof("getting image %s of user %s", imageID, userID)

	if image, err := interactor.storageDB.getImage(userID, imageID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting image on storage database %s", err)
		return nil, err
	} else {
		return image, nil
	}
}

// getImageRaw ...
func (interactor *interactor) getImageRaw(userID string, imageID string) ([]byte, error) {
	log.WithFields(map[string]interface{}{"method": "getImage"})
	log.Infof("getting rawImage %s of user %s", imageID, userID)

	if interactor.config.Dropbox.Enabled {
		path := fmt.Sprintf("/users/%s/images/%s", userID, imageID)
		if rawImage, err := interactor.storageDropbox.download(path); err != nil {
			log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
				Errorf("error getting rawImage on storage dropbox %s", err)
			return nil, err
		} else {
			return rawImage, nil
		}
	} else {
		if image, err := interactor.storageDB.getImage(userID, imageID); err != nil {
			log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
				Errorf("error getting image on storage database %s", err)
			return nil, err
		} else {
			return image.RawImage, nil
		}
	}
}

// createImage ...
func (interactor *interactor) createImage(newImage *image) (*image, error) {
	log.WithFields(map[string]interface{}{"method": "createImage"})

	log.Info("creating image")
	newImage.ImageID = genUI()

	if image, err := interactor.storageDB.createImage(newImage); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error creating image on storage database %s", err)
		return nil, err
	} else {
		if interactor.config.Dropbox.Enabled {
			path := fmt.Sprintf("/users/%s/images/%s", newImage.UserID, newImage.ImageID)
			if err := interactor.storageDropbox.upload(path, newImage.RawImage); err != nil {
				log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
					Errorf("error creating image on storage dropbox %s", err)
				return nil, err
			}
		}

		return image, nil
	}
}

// updateImage ...
func (interactor *interactor) updateImage(updImage *image) (*image, error) {
	log.WithFields(map[string]interface{}{"method": "updateImage"})
	log.Infof("updating image %s of user %s", updImage.UserID, updImage.UserID)
	if image, err := interactor.storageDB.updateImage(updImage); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error updating image on storage database %s", err)
		return nil, err
	} else {
		if interactor.config.Dropbox.Enabled {
			path := fmt.Sprintf("/users/%s/images/%s", updImage.UserID, updImage.ImageID)
			if err := interactor.storageDropbox.upload(path, updImage.RawImage); err != nil {
				log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
					Errorf("error updating image on storage dropbox %s", err)
				return nil, err
			}
			return image, nil
		}
	}
	return nil, nil
}

// deleteImage ...
func (interactor *interactor) deleteImage(userID string, imageID string) error {
	log.WithFields(map[string]interface{}{"method": "deleteImage"})
	log.Infof("deleting image %s of user %s", imageID, userID)
	if err := interactor.storageDB.deleteImage(userID, imageID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error deleting image on storage database %s", err)
		return err
	} else {
		if interactor.config.Dropbox.Enabled {
			path := fmt.Sprintf("/users/%s/images/%s", userID, imageID)
			if err := interactor.storageDropbox.delete(path); err != nil {
				log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
					Errorf("error deleting image on storage dropbox %s", err)
				return err
			}
		}
	}
	return nil
}

// getCategories ...
func (interactor *interactor) getCategories(userID string) ([]*category, error) {
	log.WithFields(map[string]interface{}{"method": "getCategories"})
	log.Infof("getting categories of user %s", userID)
	if categories, err := interactor.storageDB.getCategories(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting categories on storage database %s", err)
		return nil, err
	} else {
		return categories, nil
	}
}

// getCategory ...
func (interactor *interactor) getCategory(userID string, categoryID string) (*category, error) {
	log.WithFields(map[string]interface{}{"method": "getCategory"})
	log.Infof("getting category %s of user %s", categoryID, userID)
	if category, err := interactor.storageDB.getCategory(userID, categoryID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting category on storage database %s", err)
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
		category.CategoryID = genUI()
	}

	if categories, err := interactor.storageDB.createCategories(newCategories); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error creating categories on storage database %s", err)
		return nil, err
	} else {
		return categories, nil
	}
}

// updateCategory ...
func (interactor *interactor) updateCategory(updCategory *category) (*category, error) {
	log.WithFields(map[string]interface{}{"method": "updateCategory"})
	log.Infof("updating category %s of user %s", updCategory.UserID, updCategory.UserID)
	if category, err := interactor.storageDB.updateCategory(updCategory); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error updating category on storage database %s", err)
		return nil, err
	} else {
		return category, nil
	}
}

// deleteCategory ...
func (interactor *interactor) deleteCategory(userID string, categoryID string) error {
	log.WithFields(map[string]interface{}{"method": "deleteCategory"})
	log.Infof("deleting category %s of user %s", categoryID, userID)
	if err := interactor.storageDB.deleteCategory(userID, categoryID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error deleting category on storage database %s", err)
		return err
	}
	return nil
}

// getTransaction ...
func (interactor *interactor) getTransactions(userID string) ([]*transaction, error) {
	log.WithFields(map[string]interface{}{"method": "getTransactions"})
	log.Infof("getting transactions of user %s", userID)
	if transaction, err := interactor.storageDB.getTransactions(userID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting transactions on storage database %s", err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// getTransaction ...
func (interactor *interactor) getTransaction(userID string, walletID string, transactionID string) (*transaction, error) {
	log.WithFields(map[string]interface{}{"method": "getTransaction"})
	log.Infof("getting transaction %s of user %s on wallet %s", transactionID, userID, walletID)
	if transaction, err := interactor.storageDB.getTransaction(userID, walletID, transactionID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error getting transaction on storage database %s", err)
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
		transaction.TransactionID = genUI()
	}

	if transactions, err := interactor.storageDB.createTransactions(newTransactions); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error creating transactions on storage database %s", err)
		return nil, err
	} else {
		return transactions, nil
	}
}

// updateTransaction ...
func (interactor *interactor) updateTransaction(updTransaction *transaction) (*transaction, error) {
	log.WithFields(map[string]interface{}{"method": "updateTransaction"})
	log.Infof("updating transaction %s of user %s", updTransaction.UserID, updTransaction.UserID)
	if transaction, err := interactor.storageDB.updateTransaction(updTransaction); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error updating transaction on storage database %s", err)
		return nil, err
	} else {
		return transaction, nil
	}
}

// deleteTransaction ...
func (interactor *interactor) deleteTransaction(userID string, walletID string, transactionID string) error {
	log.WithFields(map[string]interface{}{"method": "deleteTransaction"})
	log.Infof("deleting transaction %s of user %s", transactionID, userID)
	if err := interactor.storageDB.deleteTransaction(userID, walletID, transactionID); err != nil {
		log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err}).
			Errorf("error deleting transaction on storage database %s", err)
		return err
	}
	return nil
}
