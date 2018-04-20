package gomoney

import (
	"net/http"

	"time"

	"github.com/joaosoft/go-manager/service"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"gopkg.in/validator.v2"
)

// apiWeb ...
type apiWeb struct {
	host       string
	interactor *interactor
}

// USERS
type createUserRequest struct {
	Body struct {
		Name        string `json:"name" validate:"nonzero"`
		Email       string `json:"email" validate:"nonzero"`
		Password    string `json:"password" validate:"nonzero"`
		Description string `json:"description"`
	}
}

type updateUserRequest struct {
	UserID string `json:"user_id" validate:"nonzero"`
	Body   struct {
		Name        string `json:"name" validate:"nonzero"`
		Email       string `json:"email" validate:"nonzero"`
		Password    string `json:"password" validate:"nonzero"`
		Description string `json:"description"`
	}
}

type userResponse struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

// WALLETS
type createWalletsRequest struct {
	UserID string              `json:"user_id" validate:"nonzero"`
	Body   []walletItemRequest `json:"wallets" validate:"min=1"`
}

type updateWalletRequest struct {
	UserID   string `json:"user_id" validate:"nonzero"`
	WalletID string `json:"wallet_id" validate:"nonzero"`
	Body     walletItemRequest
}

type walletItemRequest struct {
	Name        string `json:"name" validate:"nonzero"`
	Description string `json:"description"`
	Password    string `json:"password"`
}

type walletResponse struct {
	WalletID    string `json:"wallet_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Password    string `json:"password"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

// IMAGES
type createImagesRequest struct {
	UserID string `json:"user_id" validate:"nonzero"`
	Body   []struct {
		Name        string `json:"name" validate:"nonzero"`
		Description string `json:"description" validate:"nonzero"`
		Url         string `json:"url"`
		RawImage    []byte `json:"raw_image"`
	} `json:"images" validate:"min=1"`
}

type updateImageRequest struct {
	ImageID string `json:"image_id" validate:"nonzero"`
	UserID  string `json:"user_id" validate:"nonzero"`
	Body    struct {
		Name        string `json:"name" validate:"nonzero"`
		Description string `json:"description"`
		Url         string `json:"url"`
		RawImage    []byte `json:"raw_image"`
	} `json:"body"`
}

type imageResponse struct {
	ImageID     string `json:"image_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	RawImage    []byte `json:"raw_image"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

// CATEGORIES
type createCategoriesRequest struct {
	UserID string                `json:"user_id" validate:"nonzero"`
	Body   []categoryItemRequest `json:"categories" validate:"min=1"`
}

type updateCategoryRequest struct {
	UserID     string `json:"user_id" validate:"nonzero"`
	CategoryID string `json:"category_id" validate:"nonzero"`
	Body       categoryItemRequest
}

type categoryItemRequest struct {
	Name        string `json:"name" validate:"nonzero"`
	Description string `json:"description"`
	ImageID     string `json:"image_id"`
}

type categoryResponse struct {
	CategoryID  string `json:"category_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageID     string `json:"image_id"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

// TRANSACTIONS
type createTransactionsRequest struct {
	UserID   string                   `json:"user_id" validate:"nonzero"`
	WalletID string                   `json:"wallet_id" validate:"nonzero"`
	Body     []transactionItemRequest `json:"transactions" validate:"min=1"`
}

type updateTransactionRequest struct {
	UserID        string `json:"user_id" validate:"nonzero"`
	WalletID      string `json:"wallet_id" validate:"nonzero"`
	TransactionID string `json:"transaction_id" validate:"nonzero"`
	Body          transactionItemRequest
}

type transactionItemRequest struct {
	CategoryID  string `json:"category_id" validate:"nonzero"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Date        string `json:"date" validate:"nonzero"`
}

type transactionResponse struct {
	UserID        string `json:"user_id"`
	WalletID      string `json:"wallet_id"`
	TransactionID string `json:"transaction_id"`
	CategoryID    string `json:"category_id"`
	Price         string `json:"price"`
	Description   string `json:"description"`
	Date          string `json:"date"`
	UpdatedAt     string `json:"updated_at"`
	CreatedAt     string `json:"created_at"`
}

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// newApiWeb ...
func newApiWeb(host string, interactor *interactor) *apiWeb {
	webApi := &apiWeb{
		host:       host,
		interactor: interactor,
	}

	return webApi
}

func (api *apiWeb) new() gomanager.IWeb {
	web := gomanager.NewSimpleWebEcho(api.host)

	// users
	web.AddRoute("GET", "/users", api.getUsersHandler)
	web.AddRoute("GET", "/users/:user_id", api.getUserHandler)
	web.AddRoute("POST", "/users", api.createUserHandler)
	web.AddRoute("PUT", "/users/:user_id", api.updateUserHandler)
	web.AddRoute("DELETE", "/users/:user_id", api.deleteUserHandler)

	// wallets
	web.AddRoute("GET", "/users/:user_id/wallets", api.getWalletsHandler)
	web.AddRoute("GET", "/users/:user_id/wallets/:wallet_id", api.getWalletHandler)
	web.AddRoute("POST", "/users/:user_id/wallets", api.createWalletsHandler)
	web.AddRoute("PUT", "/users/:user_id/wallets/:wallet_id", api.updateWalletHandler)
	web.AddRoute("DELETE", "/users/:user_id/wallets/:wallet_id", api.deleteWalletHandler)

	// images
	web.AddRoute("GET", "/users/:user_id/images", api.getImagesHandler)
	web.AddRoute("GET", "/users/:user_id/images/:image_id", api.getImageHandler)
	web.AddRoute("POST", "/users/:user_id/images", api.createImagesHandler)
	web.AddRoute("PUT", "/users/:user_id/images/:image_id", api.updateImageHandler)
	web.AddRoute("DELETE", "/users/:user_id/images/:image_id", api.deleteImageHandler)

	// categories
	web.AddRoute("GET", "/users/:user_id/categories", api.getCategoriesHandler)
	web.AddRoute("GET", "/users/:user_id/categories/:category_id", api.getCategoryHandler)
	web.AddRoute("POST", "/users/:user_id/categories", api.createCategoriesHandler)
	web.AddRoute("PUT", "/users/:user_id/categories/:category_id", api.updateCategoryHandler)
	web.AddRoute("DELETE", "/users/:user_id/categories/:category_id", api.deleteCategoryHandler)

	// transactions
	web.AddRoute("GET", "/users/:user_id/wallets/:wallet_id/transactions", api.getTransactionsHandler)
	web.AddRoute("GET", "/users/:user_id/wallets/:wallet_id/transactions/:transaction_id", api.getTransactionHandler)
	web.AddRoute("POST", "/users/:user_id/wallets/:wallet_id/transactions", api.createTransactionsHandler)
	web.AddRoute("PUT", "/users/:user_id/wallets/:wallet_id/transactions/:transaction_id", api.updateTransactionHandler)
	web.AddRoute("DELETE", "/users/:user_id/wallets/:wallet_id/transactions/:transaction_id", api.deleteTransactionHandler)

	return web
}

func (api *apiWeb) getUsersHandler(ctx echo.Context) error {
	if users, err := api.interactor.getUsers(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if users == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		usersResponse := make([]*userResponse, 0)

		for _, user := range users {
			userResponse := &userResponse{
				UserID:      user.UserID.String(),
				Name:        user.Name,
				Email:       user.Email,
				Password:    user.Password,
				Description: user.Description,
				CreatedAt:   user.CreatedAt.String(),
				UpdatedAt:   user.UpdatedAt.String(),
			}
			usersResponse = append(usersResponse, userResponse)
		}
		return ctx.JSON(http.StatusOK, usersResponse)
	}
}

func (api *apiWeb) getUserHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if user, err := api.interactor.getUser(userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if user == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			userResponse{
				UserID:      user.UserID.String(),
				Name:        user.Name,
				Email:       user.Email,
				Password:    user.Password,
				Description: user.Description,
				CreatedAt:   user.CreatedAt.String(),
				UpdatedAt:   user.UpdatedAt.String(),
			})
	}
}

func (api *apiWeb) createUserHandler(ctx echo.Context) error {
	request := createUserRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := validator.Validate(request.Body); err != nil {
		log.Errorf("error when validating body request: %s", err.Error()).ToError(&err)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if createdUser, err := api.interactor.createUser(
		&user{
			Name:        request.Body.Name,
			Email:       request.Body.Email,
			Password:    request.Body.Password,
			Description: request.Body.Description,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if createdUser == nil {
		return ctx.NoContent(http.StatusInternalServerError)
	} else {
		return ctx.JSON(http.StatusCreated, userResponse{
			UserID:      createdUser.UserID.String(),
			Name:        createdUser.Name,
			Email:       createdUser.Email,
			Password:    createdUser.Password,
			Description: createdUser.Description,
			CreatedAt:   createdUser.CreatedAt.String(),
			UpdatedAt:   createdUser.UpdatedAt.String(),
		})
	}
}

func (api *apiWeb) updateUserHandler(ctx echo.Context) error {
	request := updateUserRequest{UserID: ctx.Param("user_id")}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := validator.Validate(request.Body); err != nil {
		log.Errorf("error when validating body request: %s", err.Error()).ToError(&err)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	userID, err := uuid.FromString(request.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if updatedUser, err := api.interactor.updateUser(
		&user{
			UserID:      userID,
			Name:        request.Body.Name,
			Email:       request.Body.Email,
			Password:    request.Body.Password,
			Description: request.Body.Description,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if updatedUser == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, userResponse{
			UserID:      updatedUser.UserID.String(),
			Name:        updatedUser.Name,
			Email:       updatedUser.Email,
			Password:    updatedUser.Password,
			Description: updatedUser.Description,
			CreatedAt:   updatedUser.CreatedAt.String(),
			UpdatedAt:   updatedUser.UpdatedAt.String(),
		})
	}
}

func (api *apiWeb) deleteUserHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := api.interactor.deleteUser(userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (api *apiWeb) getWalletsHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if wallets, err := api.interactor.getWallets(userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if wallets == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		walletsResponse := make([]*walletResponse, 0)
		for _, wallet := range wallets {
			walletResponse := &walletResponse{
				WalletID:    wallet.WalletID.String(),
				UserID:      wallet.UserID.String(),
				Name:        wallet.Name,
				Description: wallet.Description,
				Password:    wallet.Password,
				CreatedAt:   wallet.CreatedAt.String(),
				UpdatedAt:   wallet.UpdatedAt.String(),
			}
			walletsResponse = append(walletsResponse, walletResponse)
		}
		return ctx.JSON(http.StatusOK, walletsResponse)
	}
}

func (api *apiWeb) getWalletHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	walletID, err := uuid.FromString(ctx.Param("wallet_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if wallet, err := api.interactor.getWallet(userID, walletID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if wallet == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			walletResponse{
				WalletID:    wallet.WalletID.String(),
				Name:        wallet.Name,
				Description: wallet.Description,
				Password:    wallet.Password,
				CreatedAt:   wallet.CreatedAt.String(),
				UpdatedAt:   wallet.UpdatedAt.String(),
			})
	}
}

func (api *apiWeb) createWalletsHandler(ctx echo.Context) error {
	request := createWalletsRequest{
		UserID: ctx.Param("user_id"),
	}
	wallets := make([]*wallet, 0)

	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for _, item := range request.Body {
		if err := validator.Validate(item); err != nil {
			log.Errorf("error when validating body request: %s", err.Error()).ToError(&err)
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		}
	}

	userID, err := uuid.FromString(request.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for _, item := range request.Body {
		wallets = append(wallets, &wallet{
			UserID:      userID,
			Name:        item.Name,
			Description: item.Description,
			Password:    item.Password,
		})
	}

	if createdWallets, err := api.interactor.createWallets(wallets); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		walletsResponse := make([]*walletResponse, 0)

		for _, createdWallet := range createdWallets {
			walletResponse := &walletResponse{
				WalletID:    createdWallet.WalletID.String(),
				UserID:      createdWallet.UserID.String(),
				Name:        createdWallet.Name,
				Description: createdWallet.Description,
				Password:    createdWallet.Password,
				CreatedAt:   createdWallet.CreatedAt.String(),
				UpdatedAt:   createdWallet.UpdatedAt.String(),
			}
			walletsResponse = append(walletsResponse, walletResponse)
		}
		return ctx.JSON(http.StatusCreated, walletsResponse)
	}
}

func (api *apiWeb) updateWalletHandler(ctx echo.Context) error {
	request := updateWalletRequest{
		UserID:   ctx.Param("user_id"),
		WalletID: ctx.Param("wallet_id"),
	}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := validator.Validate(request.Body); err != nil {
		log.Errorf("error when validating body request: %s", err.Error()).ToError(&err)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	walletID, err := uuid.FromString(ctx.Param("wallet_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if updatedWallet, err := api.interactor.updateWallet(
		&wallet{
			UserID:      userID,
			WalletID:    walletID,
			Name:        request.Body.Name,
			Description: request.Body.Description,
			Password:    request.Body.Password,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if updatedWallet == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, walletResponse{
			WalletID:    updatedWallet.WalletID.String(),
			UserID:      updatedWallet.UserID.String(),
			Name:        updatedWallet.Name,
			Description: updatedWallet.Description,
			Password:    updatedWallet.Password,
			CreatedAt:   updatedWallet.CreatedAt.String(),
			UpdatedAt:   updatedWallet.UpdatedAt.String(),
		})
	}
}

func (api *apiWeb) deleteWalletHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	walletID, err := uuid.FromString(ctx.Param("wallet_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := api.interactor.deleteWallet(userID, walletID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (api *apiWeb) getImagesHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if images, err := api.interactor.getImages(userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if images == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		imagesResponse := make([]*imageResponse, 0)

		for _, image := range images {
			imageResponse := &imageResponse{
				UserID:      image.UserID.String(),
				ImageID:     image.ImageID.String(),
				Name:        image.Name,
				Description: image.Description,
				Url:         image.Url,
				RawImage:    image.RawImage,
				CreatedAt:   image.CreatedAt.String(),
				UpdatedAt:   image.UpdatedAt.String(),
			}
			imagesResponse = append(imagesResponse, imageResponse)
		}
		return ctx.JSON(http.StatusOK, imagesResponse)
	}
}

func (api *apiWeb) getImageHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}
	imageID, err := uuid.FromString(ctx.Param("image_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if image, err := api.interactor.getImage(userID, imageID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if image == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			imageResponse{
				UserID:      image.UserID.String(),
				ImageID:     image.ImageID.String(),
				Name:        image.Name,
				Description: image.Description,
				Url:         image.Url,
				RawImage:    image.RawImage,
				CreatedAt:   image.CreatedAt.String(),
				UpdatedAt:   image.UpdatedAt.String(),
			})
	}
}

func (api *apiWeb) createImagesHandler(ctx echo.Context) error {
	request := createImagesRequest{
		UserID: ctx.Param("user_id"),
	}
	images := make([]*image, 0)

	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for _, item := range request.Body {
		if err := validator.Validate(item); err != nil {
			log.Errorf("error when validating body request: %s", err.Error()).ToError(&err)
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		}
	}

	userID, err := uuid.FromString(request.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for _, item := range request.Body {
		images = append(images, &image{
			UserID:      userID,
			Name:        item.Name,
			Description: item.Description,
			Url:         item.Url,
			RawImage:    item.RawImage,
		})
	}

	if createdImages, err := api.interactor.createImages(images); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		imagesResponse := make([]*imageResponse, 0)

		for _, createdImage := range createdImages {
			imageResponse := &imageResponse{
				ImageID:     createdImage.ImageID.String(),
				UserID:      createdImage.UserID.String(),
				Name:        createdImage.Name,
				Description: createdImage.Description,
				Url:         createdImage.Url,
				RawImage:    createdImage.RawImage,
				CreatedAt:   createdImage.CreatedAt.String(),
				UpdatedAt:   createdImage.UpdatedAt.String(),
			}
			imagesResponse = append(imagesResponse, imageResponse)
		}
		return ctx.JSON(http.StatusCreated, imagesResponse)
	}
}

func (api *apiWeb) updateImageHandler(ctx echo.Context) error {
	request := updateImageRequest{
		UserID:  ctx.Param("user_id"),
		ImageID: ctx.Param("image_id"),
	}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := validator.Validate(request.Body); err != nil {
		log.Errorf("error when validating body request: %s", err.Error()).ToError(&err).ToError(&err)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	imageID, err := uuid.FromString(request.ImageID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	userID, err := uuid.FromString(request.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if updatedImage, err := api.interactor.updateImage(
		&image{
			ImageID:     imageID,
			UserID:      userID,
			Name:        request.Body.Name,
			Description: request.Body.Description,
			Url:         request.Body.Url,
			RawImage:    request.Body.RawImage,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if updatedImage == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, imageResponse{
			ImageID:     updatedImage.ImageID.String(),
			UserID:      updatedImage.UserID.String(),
			Name:        updatedImage.Name,
			Description: updatedImage.Description,
			Url:         updatedImage.Url,
			RawImage:    updatedImage.RawImage,
			CreatedAt:   updatedImage.CreatedAt.String(),
			UpdatedAt:   updatedImage.UpdatedAt.String(),
		})
	}
}

func (api *apiWeb) deleteImageHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	imageID, err := uuid.FromString(ctx.Param("image_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := api.interactor.deleteImage(userID, imageID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (api *apiWeb) getCategoriesHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if categories, err := api.interactor.getCategories(userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if categories == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		categoriesResponse := make([]*categoryResponse, 0)

		for _, category := range categories {
			categoryResponse := &categoryResponse{
				CategoryID:  category.CategoryID.String(),
				UserID: category.UserID.String(),
				Name:        category.Name,
				Description: category.Description,
				ImageID:     category.ImageID.String(),
				CreatedAt:   category.CreatedAt.String(),
				UpdatedAt:   category.UpdatedAt.String(),
			}
			categoriesResponse = append(categoriesResponse, categoryResponse)
		}
		return ctx.JSON(http.StatusOK, categoriesResponse)
	}
}

func (api *apiWeb) getCategoryHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	categoryID, err := uuid.FromString(ctx.Param("category_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if category, err := api.interactor.getCategory(userID, categoryID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if category == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			categoryResponse{
				CategoryID:  category.CategoryID.String(),
				UserID:      category.UserID.String(),
				Name:        category.Name,
				Description: category.Description,
				ImageID:     category.ImageID.String(),
				CreatedAt:   category.CreatedAt.String(),
				UpdatedAt:   category.UpdatedAt.String(),
			})
	}
}

func (api *apiWeb) createCategoriesHandler(ctx echo.Context) error {
	request := createCategoriesRequest{
		UserID: ctx.Param("user_id"),
	}
	categories := make([]*category, 0)

	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for _, item := range request.Body {
		if err := validator.Validate(item); err != nil {
			log.Errorf("error when validating body request: %s", err.Error()).ToError(&err)
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		}
	}

	userID, err := uuid.FromString(request.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for _, item := range request.Body {
		imageID, err := uuid.FromString(item.ImageID)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		}

		categories = append(categories, &category{
			UserID:      userID,
			Name:        item.Name,
			Description: item.Description,
			ImageID:     imageID,
		})
	}

	if createdCategories, err := api.interactor.createCategories(categories); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		categoriesResponse := make([]*categoryResponse, 0)

		for _, createdCategory := range createdCategories {
			categoryResponse := &categoryResponse{
				CategoryID:  createdCategory.CategoryID.String(),
				UserID:      createdCategory.UserID.String(),
				Name:        createdCategory.Name,
				Description: createdCategory.Description,
				ImageID:     createdCategory.ImageID.String(),
				CreatedAt:   createdCategory.CreatedAt.String(),
				UpdatedAt:   createdCategory.UpdatedAt.String(),
			}
			categoriesResponse = append(categoriesResponse, categoryResponse)
		}
		return ctx.JSON(http.StatusCreated, categoriesResponse)
	}
}

func (api *apiWeb) updateCategoryHandler(ctx echo.Context) error {
	request := updateCategoryRequest{
		UserID:     ctx.Param("user_id"),
		CategoryID: ctx.Param("category_id"),
	}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := validator.Validate(request.Body); err != nil {
		log.Errorf("error when validating body request: %s", err.Error()).ToError(&err)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	categoryID, err := uuid.FromString(ctx.Param("category_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if updatedCategory, err := api.interactor.updateCategory(
		&category{
			UserID:      userID,
			CategoryID:  categoryID,
			Name:        request.Body.Name,
			Description: request.Body.Description,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if updatedCategory == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, categoryResponse{
			CategoryID:  updatedCategory.CategoryID.String(),
			UserID:      updatedCategory.UserID.String(),
			Name:        updatedCategory.Name,
			Description: updatedCategory.Description,
			ImageID:     updatedCategory.ImageID.String(),
			CreatedAt:   updatedCategory.CreatedAt.String(),
			UpdatedAt:   updatedCategory.UpdatedAt.String(),
		})
	}
}

func (api *apiWeb) deleteCategoryHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	categoryID, err := uuid.FromString(ctx.Param("category_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := api.interactor.deleteCategory(userID, categoryID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (api *apiWeb) getTransactionsHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if transactions, err := api.interactor.getTransactions(userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if transactions == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		transactionsResponse := make([]*transactionResponse, 0)

		for _, transaction := range transactions {
			transactionResponse := &transactionResponse{
				TransactionID: transaction.TransactionID.String(),
				UserID:        transaction.UserID.String(),
				WalletID:      transaction.WalletID.String(),
				CategoryID:    transaction.CategoryID.String(),
				Price:         transaction.Price.String(),
				Description:   transaction.Description,
				Date:          transaction.Date.String(),
				CreatedAt:     transaction.CreatedAt.String(),
				UpdatedAt:     transaction.UpdatedAt.String(),
			}
			transactionsResponse = append(transactionsResponse, transactionResponse)
		}
		return ctx.JSON(http.StatusOK, transactionsResponse)
	}
}

func (api *apiWeb) getTransactionHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	walletID, err := uuid.FromString(ctx.Param("wallet_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	transactionID, err := uuid.FromString(ctx.Param("transaction_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if transaction, err := api.interactor.getTransaction(userID, walletID, transactionID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if transaction == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			transactionResponse{
				TransactionID: transaction.TransactionID.String(),
				UserID:        transaction.UserID.String(),
				WalletID:      transaction.WalletID.String(),
				CategoryID:    transaction.CategoryID.String(),
				Price:         transaction.Price.String(),
				Description:   transaction.Description,
				Date:          transaction.Date.String(),
				CreatedAt:     transaction.CreatedAt.String(),
				UpdatedAt:     transaction.UpdatedAt.String(),
			})
	}
}

func (api *apiWeb) createTransactionsHandler(ctx echo.Context) error {
	request := createTransactionsRequest{
		UserID:   ctx.Param("user_id"),
		WalletID: ctx.Param("wallet_id"),
	}
	transactions := make([]*transaction, 0)

	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for _, item := range request.Body {
		if err := validator.Validate(item); err != nil {
			log.Errorf("error when validating body request: %s", err.Error()).ToError(&err)
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		}
	}

	userID, err := uuid.FromString(request.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	walletID, err := uuid.FromString(request.WalletID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for _, item := range request.Body {
		categoryID, err := uuid.FromString(item.CategoryID)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		}

		date, err := time.Parse(time.RFC3339, item.Date)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
		}

		price, err := decimal.NewFromString(item.Price)

		transactions = append(transactions, &transaction{
			UserID:      userID,
			WalletID:    walletID,
			CategoryID:  categoryID,
			Price:       price,
			Description: item.Description,
			Date:        date,
		})
	}

	if createdTransactions, err := api.interactor.createTransactions(transactions); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		transactionsResponse := make([]*transactionResponse, 0)

		for _, createdTransaction := range createdTransactions {
			transactionResponse := &transactionResponse{
				TransactionID: createdTransaction.TransactionID.String(),
				UserID:        createdTransaction.UserID.String(),
				WalletID:      createdTransaction.WalletID.String(),
				CategoryID:    createdTransaction.CategoryID.String(),
				Price:         createdTransaction.Price.String(),
				Description:   createdTransaction.Description,
				Date:          createdTransaction.Date.String(),
				CreatedAt:     createdTransaction.CreatedAt.String(),
				UpdatedAt:     createdTransaction.UpdatedAt.String(),
			}
			transactionsResponse = append(transactionsResponse, transactionResponse)
		}
		return ctx.JSON(http.StatusCreated, transactionsResponse)
	}
}

func (api *apiWeb) updateTransactionHandler(ctx echo.Context) error {
	request := updateTransactionRequest{
		UserID:        ctx.Param("user_id"),
		WalletID:      ctx.Param("wallet_id"),
		TransactionID: ctx.Param("transaction_id"),
	}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := validator.Validate(request.Body); err != nil {
		log.Errorf("error when validating body request: %s", err.Error()).ToError(&err)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	userID, err := uuid.FromString(request.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	transactionID, err := uuid.FromString(request.TransactionID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	categoryID, err := uuid.FromString(request.Body.CategoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	price, err := decimal.NewFromString(request.Body.Price)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	date, err := time.Parse(time.RFC3339, request.Body.Date)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if updatedTransaction, err := api.interactor.updateTransaction(
		&transaction{
			UserID:        userID,
			TransactionID: transactionID,
			CategoryID:    categoryID,
			Price:         price,
			Description:   request.Body.Description,
			Date:          date,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else if updatedTransaction == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, transactionResponse{
			TransactionID: updatedTransaction.TransactionID.String(),
			UserID:        updatedTransaction.UserID.String(),
			WalletID:      updatedTransaction.WalletID.String(),
			CategoryID:    updatedTransaction.CategoryID.String(),
			Price:         updatedTransaction.Price.String(),
			Description:   updatedTransaction.Description,
			Date:          updatedTransaction.Date.String(),
			CreatedAt:     updatedTransaction.CreatedAt.String(),
			UpdatedAt:     updatedTransaction.UpdatedAt.String(),
		})
	}
}

func (api *apiWeb) deleteTransactionHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	walletID, err := uuid.FromString(ctx.Param("wallet_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	transactionID, err := uuid.FromString(ctx.Param("transaction_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if err := api.interactor.deleteTransaction(userID, walletID, transactionID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}
