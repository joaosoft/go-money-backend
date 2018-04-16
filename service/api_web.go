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
	interactor *Interactor
}

type createUserRequest struct {
	Body struct {
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

type createTransactionRequest struct {
	Body struct {
		UserID      string `json:"user_id"`
		CategoryID  string `json:"category_id"`
		Price       string `json:"price"`
		Description string `json:"description"`
		Date        string `json:"date"`
	}
}

type transactionResponse struct {
	UserID        string `json:"user_id"`
	TransactionID string `json:"transaction_id"`
	CategoryID    string `json:"category_id"`
	Price         string `json:"price"`
	Description   string `json:"description"`
	Date          string `json:"date"`
	UpdatedAt     string `json:"updated_at"`
	CreatedAt     string `json:"created_at"`
}

// newApiWeb ...
func newApiWeb(host string, interactor *Interactor) *apiWeb {
	webApi := &apiWeb{
		host:       host,
		interactor: interactor,
	}

	return webApi
}

func (api *apiWeb) new() gomanager.IWeb {
	web := gomanager.NewSimpleWebEcho(api.host)

	// user
	web.AddRoute("GET", "/users/:user_id", api.getUserHandler)
	web.AddRoute("POST", "/users", api.createUserHandler)
	web.AddRoute("PUT", "/users/:user_id", api.updateUserHandler)
	web.AddRoute("DELETE", "/users/:user_id", api.deleteUserHandler)

	// transactions
	web.AddRoute("GET", "/transactions/:transaction_id", api.getTransactionHandler)
	web.AddRoute("POST", "/transactions", api.createTransactionHandler)
	web.AddRoute("PUT", "/transactions/:transaction_id", api.updateTransactionHandler)
	web.AddRoute("DELETE", "/transactions/:transaction_id", api.deleteTransactionHandler)

	return web
}

func (api *apiWeb) getUserHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if user, err := api.interactor.GetUser(userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
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
				UpdatedAt:   user.UpdatedAt.String(),
				CreatedAt:   user.CreatedAt.String(),
			})
	}
}

func (api *apiWeb) createUserHandler(ctx echo.Context) error {
	request := createUserRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validator.Validate(request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if createdUser, err := api.interactor.CreateUser(
		&User{
			Name:        request.Body.Name,
			Email:       request.Body.Email,
			Password:    request.Body.Password,
			Description: request.Body.Description,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else if createdUser == nil {
		return ctx.NoContent(http.StatusInternalServerError)
	} else {
		return ctx.JSON(http.StatusCreated, createdUser)
	}
}

func (api *apiWeb) updateUserHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	request := createUserRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validator.Validate(request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if updatedUser, err := api.interactor.UpdateUser(
		&User{
			UserID:      userID,
			Name:        request.Body.Name,
			Email:       request.Body.Email,
			Password:    request.Body.Password,
			Description: request.Body.Description,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else if updatedUser == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, updatedUser)
	}
}

func (api *apiWeb) deleteUserHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := api.interactor.DeleteUser(userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

func (api *apiWeb) getTransactionHandler(ctx echo.Context) error {
	transactionID, err := uuid.FromString(ctx.Param("transaction_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if transaction, err := api.interactor.GetTransaction(transactionID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else if transaction == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			transactionResponse{
				TransactionID: transaction.TransactionID.String(),
				UserID:        transaction.UserID.String(),
				CategoryID:    transaction.CategoryID.String(),
				Price:         transaction.Price.String(),
				Description:   transaction.Description,
				Date:          transaction.Date.String(),
				UpdatedAt:     transaction.UpdatedAt.String(),
				CreatedAt:     transaction.CreatedAt.String(),
			})
	}
}

func (api *apiWeb) createTransactionHandler(ctx echo.Context) error {
	request := createTransactionRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validator.Validate(request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	userID, err := uuid.FromString(request.Body.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	categoryID, err := uuid.FromString(request.Body.CategoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	date, err := time.Parse(time.RFC3339, request.Body.Date)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	price, err := decimal.NewFromString(request.Body.Price)

	if createdTransaction, err := api.interactor.CreateTransaction(
		&Transaction{
			UserID:      userID,
			CategoryID:  categoryID,
			Price:       price,
			Description: request.Body.Description,
			Date:        date,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else if createdTransaction == nil {
		return ctx.NoContent(http.StatusInternalServerError)
	} else {
		return ctx.JSON(http.StatusCreated, createdTransaction)
	}
}

func (api *apiWeb) updateTransactionHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("transaction_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	request := createTransactionRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validator.Validate(request.Body); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	categoryID, err := uuid.FromString(request.Body.CategoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	date, err := time.Parse(time.RFC3339, request.Body.Date)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	price, err := decimal.NewFromString(request.Body.Price)

	if updatedTransaction, err := api.interactor.UpdateTransaction(
		&Transaction{
			UserID:      userID,
			CategoryID:  categoryID,
			Price:       price,
			Description: request.Body.Description,
			Date:        date,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else if updatedTransaction == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, updatedTransaction)
	}
}

func (api *apiWeb) deleteTransactionHandler(ctx echo.Context) error {
	transactionID, err := uuid.FromString(ctx.Param("transaction_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := api.interactor.DeleteTransaction(transactionID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}
