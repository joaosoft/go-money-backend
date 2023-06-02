// Go Money API
//
// Go Money is a general day-to-day expenses manager.
//
//	Schemes: http
//	BasePath: /api/1
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package gomoney

import (
	"fmt"
	"github.com/joaosoft/manager"
	"net/http"
	"strings"
	"time"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/validator"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

const (
	tokenName      = "AccessToken"
	authentication = "Bearer"
	session_key    = "Authorization"
)

// apiWeb ...
type apiWeb struct {
	host       string
	auth       echo.MiddlewareFunc
	client     manager.IWeb
	interactor *interactor
}

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Cause   string `json:"cause,omitempty"`
}

func (m *Money) newApiWeb(host string, interactor *interactor) *apiWeb {
	webApi := &apiWeb{
		host:       host,
		interactor: interactor,
		client:     m.pm.NewSimpleWebEcho(host),
	}

	webApi.registerRoutes()

	return webApi
}

func (api *apiWeb) registerRoutes() error {
	api.client = manager.NewSimpleWebEcho(api.host)
	api.auth = api.authenticate()

	api.registerRoutesForUsers()
	api.registerRoutesForSessions()
	api.registerRoutesForWallets()
	api.registerRoutesForCategories()
	api.registerRoutesForImages()
	api.registerRoutesForTransactions()

	return nil
}

type authenticateRequest struct {
	UserID string `json:"user_id" validate:"ui"`
}

func (api *apiWeb) authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			request := authenticateRequest{UserID: ctx.Param("user_id")}
			if err := validator.Validate(request); err != nil {
				newErr := errors.NewError(err)
				log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
					Error("error when validating body request").ToErrorData(newErr)
				return ctx.JSON(http.StatusNetworkAuthenticationRequired, errorResponse{Code: http.StatusNetworkAuthenticationRequired, Message: newErr.Error(), Cause: newErr.Cause()})
			}

			sessionKeyValue := ctx.Request().Header.Get(session_key)
			sessionKeyValue = strings.Replace(sessionKeyValue, fmt.Sprintf("%s ", authentication), "", 1)

			token, err := jwt.Parse(sessionKeyValue, func(token *jwt.Token) (interface{}, error) {
				if session, err := api.interactor.getSession(request.UserID, sessionKeyValue); err != nil {
					log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
						Error("error getting session").ToErrorData(err)
					return nil, err
				} else if session == nil {
					var err error
					log.WithFields(map[string]interface{}{"error": err.Error()}).
						Error("unexisting session").ToError(&err)
					return nil, err
				} else {
					return []byte(session.Original), nil
				}
			})

			if err == nil && token.Valid {
				log.Infof("valid token %s", sessionKeyValue)
				return next(ctx)
			} else {
				log.Infof("invalid token %s with error %s", sessionKeyValue, err)
				return ctx.NoContent(http.StatusNetworkAuthenticationRequired)
			}

			return next(ctx)
		}
	}
}

type getUserRequest struct {
	UserID string `json:"user_id" validate:"ui"`
}

type createUserRequest struct {
	Body struct {
		Name        string `json:"name" validate:"nonzero"`
		Email       string `json:"email" validate:"nonzero"`
		Password    string `json:"password" validate:"nonzero"`
		Description string `json:"description"`
	}
}

type updateUserRequest struct {
	UserID string `json:"user_id" validate:"ui"`
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
	Description string `json:"description,omitempty"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

type deleteUserRequest struct {
	UserID string `json:"user_id" validate:"ui"`
}

func (api *apiWeb) registerRoutesForUsers() error {
	api.client.AddRoute(http.MethodGet, "/api/1/users", api.getUsersHandler, api.auth)
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id", api.getUserHandler, api.auth)
	api.client.AddRoute(http.MethodPost, "/api/1/users", api.createUserHandler)
	api.client.AddRoute(http.MethodPut, "/api/1/users/:user_id", api.updateUserHandler, api.auth)
	api.client.AddRoute(http.MethodDelete, "/api/1/users/:user_id", api.deleteUserHandler, api.auth)

	return nil
}

// swagger:route GET /api/1/users users
//
// Gets all users.
//
// This api gets all users.
//
//	    Consumes:
//	    - application/json
//
//	    Produces:
//	    - application/json
//
//	    Schemes: http
//
//	    Responses:
//	      200: userResponse
//	      404:
//			 500:
func (api *apiWeb) getUsersHandler(ctx echo.Context) error {
	if users, err := api.interactor.getUsers(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if users == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		usersResponse := make([]*userResponse, 0)

		for _, user := range users {
			userResponse := &userResponse{
				UserID:      user.UserID,
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

// swagger:route GET /api/1/users/{user_id} user getUserRequest
//
// Gets a user.
//
// This api gets a user.
//
//	    Consumes:
//	    - application/json
//
//	    Produces:
//	    - application/json
//
//	    Schemes: http
//
//	    Responses:
//	      200: []userResponse
//			 400:
//	      404:
//			 500:
func (api *apiWeb) getUserHandler(ctx echo.Context) error {
	request := getUserRequest{UserID: ctx.Param("user_id")}
	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if user, err := api.interactor.getUser(request.UserID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if user == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			userResponse{
				UserID:      user.UserID,
				Name:        user.Name,
				Email:       user.Email,
				Password:    user.Password,
				Description: user.Description,
				CreatedAt:   user.CreatedAt.String(),
				UpdatedAt:   user.UpdatedAt.String(),
			})
	}
}

// swagger:route POST /api/1/users user createUserRequest
//
// Creates a user.
//
// This api creates a user.
//
//	    Consumes:
//	    - application/json
//
//	    Produces:
//	    - application/json
//
//	    Schemes: http
//
//	    Responses:
//	      201: userResponse
//			 400:
//			 500:
func (api *apiWeb) createUserHandler(ctx echo.Context) error {
	request := createUserRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := validator.Validate(request.Body); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if createdUser, err := api.interactor.createUser(
		&user{
			Name:        request.Body.Name,
			Email:       request.Body.Email,
			Password:    request.Body.Password,
			Description: request.Body.Description,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if createdUser == nil {
		return ctx.NoContent(http.StatusInternalServerError)
	} else {
		return ctx.JSON(http.StatusCreated, userResponse{
			UserID:      createdUser.UserID,
			Name:        createdUser.Name,
			Email:       createdUser.Email,
			Password:    createdUser.Password,
			Description: createdUser.Description,
			CreatedAt:   createdUser.CreatedAt.String(),
			UpdatedAt:   createdUser.UpdatedAt.String(),
		})
	}
}

// swagger:route PUT /api/1/users/{user_id} user createUserRequest
//
// Updates a user.
//
// This api updates a user.
//
//	    Consumes:
//	    - application/json
//
//	    Produces:
//	    - application/json
//
//	    Schemes: http
//
//	    Responses:
//	      200: userResponse
//			 400:
//	      404:
//			 500:
func (api *apiWeb) updateUserHandler(ctx echo.Context) error {
	request := updateUserRequest{UserID: ctx.Param("user_id")}
	if err := ctx.Bind(&request.Body); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if updatedUser, err := api.interactor.updateUser(
		&user{
			UserID:      request.UserID,
			Name:        request.Body.Name,
			Email:       request.Body.Email,
			Password:    request.Body.Password,
			Description: request.Body.Description,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if updatedUser == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK, userResponse{
			UserID:      updatedUser.UserID,
			Name:        updatedUser.Name,
			Email:       updatedUser.Email,
			Password:    updatedUser.Password,
			Description: updatedUser.Description,
			CreatedAt:   updatedUser.CreatedAt.String(),
			UpdatedAt:   updatedUser.UpdatedAt.String(),
		})
	}
}

// swagger:route DELETE /api/1/users/{user_id} user deleteUserRequest
//
// Deletes a user.
//
// This api deletes a user.
//
//	    Consumes:
//	    - application/json
//
//	    Produces:
//	    - application/json
//
//	    Schemes: http
//
//	    Responses:
//	      200:
//			 400:
//			 500:
func (api *apiWeb) deleteUserHandler(ctx echo.Context) error {
	request := deleteUserRequest{UserID: ctx.Param("user_id")}
	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := api.interactor.deleteUser(request.UserID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

type createSessionRequest struct {
	Body struct {
		Email       string `json:"email" validate:"nonzero"`
		Password    string `json:"password" validate:"nonzero"`
		Description string `json:"description"`
	}
}

type deleteSessionRequest struct {
	Email string `json:"email" validate:"nonzero"`
}

type deleteSessionsRequest struct {
	UserID string `json:"user_id" validate:"ui"`
}

type sessionResponse struct {
	User        sessionUserResponse `json:"user"`
	SessionID   string              `json:"session_id"`
	Token       string              `json:"token"`
	Description string              `json:"description"`
	UpdatedAt   string              `json:"updated_at"`
	CreatedAt   string              `json:"created_at"`
}

type sessionUserResponse struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

func (api *apiWeb) registerRoutesForSessions() error {
	api.client.AddRoute(http.MethodPost, "/api/1/sessions", api.createSessionHandler)
	api.client.AddRoute(http.MethodDelete, "/api/1/users/:email/session", api.deleteSessionHandler, api.auth)
	api.client.AddRoute(http.MethodDelete, "/api/1/users/:user_id/sessions", api.deleteSessionsHandler, api.auth)

	return nil
}

func (api *apiWeb) createSessionHandler(ctx echo.Context) error {
	request := createSessionRequest{}
	if err := ctx.Bind(&request.Body); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := validator.Validate(request.Body); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if user, err := api.interactor.getUserByEmail(request.Body.Email); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Errorf("error getting user by email %s", request.Body.Email).ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		passwordToken, err := generateToken(authentication, []byte(request.Body.Password))
		if err != nil {
			newErr := errors.NewError(err)
			log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
				Error("error when comparing password").ToErrorData(newErr)
			return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: newErr.Error(), Cause: newErr.Cause()})
		}

		if user.Token != passwordToken {
			newErr := errors.FromString(fmt.Sprintf("invalid password expected: %s, given: %s", user.Token, passwordToken))
			log.WithFields(map[string]interface{}{}).Error(newErr.Error())
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
		}

		if createdSession, err := api.interactor.createSession(&session{
			UserID:      user.UserID,
			Description: request.Body.Description,
		}); err != nil {
			log.WithFields(map[string]interface{}{"error": err.Error(), "cause": err.Cause()}).
				Error("error when creating session").ToErrorData(err)
			return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
		} else if createdSession == nil {
			return ctx.NoContent(http.StatusInternalServerError)
		} else {
			token := fmt.Sprintf("%s %s", authentication, createdSession.Token)
			ctx.Response().Header().Set(session_key, token)

			ctx.SetCookie(&http.Cookie{
				Name:       session_key,
				Value:      token,
				Path:       "/api/1/",
				RawExpires: "0",
			})

			return ctx.JSON(http.StatusCreated, sessionResponse{
				User: sessionUserResponse{
					UserID:      user.UserID,
					Name:        user.Name,
					Email:       user.Email,
					Description: user.Description,
					UpdatedAt:   user.UpdatedAt.String(),
					CreatedAt:   user.CreatedAt.String(),
				},
				SessionID:   createdSession.SessionID,
				Token:       createdSession.Token,
				Description: createdSession.Description,
				CreatedAt:   createdSession.CreatedAt.String(),
				UpdatedAt:   createdSession.UpdatedAt.String(),
			})
		}
	}
}

func (api *apiWeb) deleteSessionHandler(ctx echo.Context) error {
	request := deleteSessionRequest{
		Email: ctx.Param("email"),
	}

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	token := ctx.Request().Header.Get(session_key)

	if user, err := api.interactor.getUserByEmail(request.Email); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Errorf("error getting user by email %s", request.Email).ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		if err := api.interactor.deleteSession(user.UserID, token); err != nil {
			return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
		} else {
			return ctx.NoContent(http.StatusOK)
		}
	}
}

func (api *apiWeb) deleteSessionsHandler(ctx echo.Context) error {
	request := deleteSessionsRequest{
		UserID: ctx.Param("user_id"),
	}

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := api.interactor.deleteSessions(request.UserID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

type getWalletsRequest struct {
	UserID string `json:"user_id" validate:"ui"`
}

type getWalletRequest struct {
	UserID   string `json:"user_id" validate:"ui"`
	WalletID string `json:"wallet_id" validate:"ui"`
}

type createWalletsRequest struct {
	UserID string              `json:"user_id" validate:"ui"`
	Body   []walletItemRequest `json:"wallets" validate:"min=1"`
}

type updateWalletRequest struct {
	UserID   string `json:"user_id" validate:"ui"`
	WalletID string `json:"wallet_id" validate:"ui"`
	Body     walletItemRequest
}

type walletItemRequest struct {
	Name        string `json:"name" validate:"nonzero"`
	Description string `json:"description"`
	Password    string `json:"password"`
}

type deleteWalletRequest struct {
	UserID   string `json:"user_id" validate:"ui"`
	WalletID string `json:"wallet_id" validate:"ui"`
}

type walletResponse struct {
	WalletID    string `json:"wallet_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Password    string `json:"password,omitempty"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

func (api *apiWeb) registerRoutesForWallets() error {
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id/wallets", api.getWalletsHandler, api.auth)
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id/wallets/:wallet_id", api.getWalletHandler, api.auth)
	api.client.AddRoute(http.MethodPost, "/api/1/users/:user_id/wallets", api.createWalletsHandler, api.auth)
	api.client.AddRoute(http.MethodPut, "/api/1/users/:user_id/wallets/:wallet_id", api.updateWalletHandler, api.auth)
	api.client.AddRoute(http.MethodDelete, "/api/1/users/:user_id/wallets/:wallet_id", api.deleteWalletHandler, api.auth)

	return nil
}

func (api *apiWeb) getWalletsHandler(ctx echo.Context) error {
	request := getWalletsRequest{
		UserID: ctx.Param("user_id"),
	}

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if wallets, err := api.interactor.getWallets(request.UserID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if wallets == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		walletsResponse := make([]*walletResponse, 0)
		for _, wallet := range wallets {
			walletResponse := &walletResponse{
				WalletID:    wallet.WalletID,
				UserID:      wallet.UserID,
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
	request := getWalletRequest{
		UserID:   ctx.Param("user_id"),
		WalletID: ctx.Param("wallet_id"),
	}

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if wallet, err := api.interactor.getWallet(request.UserID, request.WalletID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if wallet == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			walletResponse{
				WalletID:    wallet.WalletID,
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
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	for _, item := range request.Body {
		wallets = append(wallets, &wallet{
			UserID:      request.UserID,
			Name:        item.Name,
			Description: item.Description,
			Password:    item.Password,
		})
	}

	if createdWallets, err := api.interactor.createWallets(wallets); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		walletsResponse := make([]*walletResponse, 0)

		for _, createdWallet := range createdWallets {
			walletResponse := &walletResponse{
				WalletID:    createdWallet.WalletID,
				UserID:      createdWallet.UserID,
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
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := validator.Validate(request.Body); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if updatedWallet, err := api.interactor.updateWallet(
		&wallet{
			UserID:      request.UserID,
			WalletID:    request.WalletID,
			Name:        request.Body.Name,
			Description: request.Body.Description,
			Password:    request.Body.Password,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if updatedWallet == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, walletResponse{
			WalletID:    updatedWallet.WalletID,
			UserID:      updatedWallet.UserID,
			Name:        updatedWallet.Name,
			Description: updatedWallet.Description,
			Password:    updatedWallet.Password,
			CreatedAt:   updatedWallet.CreatedAt.String(),
			UpdatedAt:   updatedWallet.UpdatedAt.String(),
		})
	}
}

func (api *apiWeb) deleteWalletHandler(ctx echo.Context) error {
	request := updateWalletRequest{
		UserID:   ctx.Param("user_id"),
		WalletID: ctx.Param("wallet_id"),
	}
	if err := ctx.Bind(&request.Body); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := api.interactor.deleteWallet(request.UserID, request.WalletID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

type getCategoriesRequest struct {
	UserID string `json:"user_id" validate:"ui"`
}

type getCategoryRequest struct {
	UserID     string `json:"user_id" validate:"ui"`
	CategoryID string `json:"category_id" validate:"ui"`
}

type createCategoriesRequest struct {
	UserID string                `json:"user_id" validate:"ui"`
	Body   []categoryItemRequest `json:"categories" validate:"min=1"`
}

type updateCategoryRequest struct {
	UserID     string `json:"user_id" validate:"ui"`
	CategoryID string `json:"category_id" validate:"ui"`
	Body       categoryItemRequest
}

type categoryItemRequest struct {
	Name        string `json:"name" validate:"nonzero"`
	Description string `json:"description"`
	ImageID     string `json:"image_id" validate:"ui"`
}

type deleteCategoryRequest struct {
	UserID     string `json:"user_id" validate:"ui"`
	CategoryID string `json:"category_id" validate:"ui"`
}

type categoryResponse struct {
	CategoryID  string `json:"category_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ImageID     string `json:"image_id"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

func (api *apiWeb) registerRoutesForCategories() error {
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id/categories", api.getCategoriesHandler, api.auth)
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id/categories/:category_id", api.getCategoryHandler, api.auth)
	api.client.AddRoute(http.MethodPost, "/api/1/users/:user_id/categories", api.createCategoriesHandler, api.auth)
	api.client.AddRoute(http.MethodPut, "/api/1/users/:user_id/categories/:category_id", api.updateCategoryHandler, api.auth)
	api.client.AddRoute(http.MethodDelete, "/api/1/users/:user_id/categories/:category_id", api.deleteCategoryHandler, api.auth)

	return nil
}

func (api *apiWeb) getCategoriesHandler(ctx echo.Context) error {
	request := getCategoriesRequest{
		UserID: ctx.Param("user_id"),
	}

	if categories, err := api.interactor.getCategories(request.UserID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if categories == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		categoriesResponse := make([]*categoryResponse, 0)

		for _, category := range categories {
			categoryResponse := &categoryResponse{
				CategoryID:  category.CategoryID,
				UserID:      category.UserID,
				Name:        category.Name,
				Description: category.Description,
				ImageID:     category.ImageID,
				CreatedAt:   category.CreatedAt.String(),
				UpdatedAt:   category.UpdatedAt.String(),
			}
			categoriesResponse = append(categoriesResponse, categoryResponse)
		}
		return ctx.JSON(http.StatusOK, categoriesResponse)
	}
}

func (api *apiWeb) getCategoryHandler(ctx echo.Context) error {
	request := getCategoryRequest{
		UserID:     ctx.Param("user_id"),
		CategoryID: ctx.Param("category_id"),
	}

	if category, err := api.interactor.getCategory(request.UserID, request.CategoryID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if category == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			categoryResponse{
				CategoryID:  category.CategoryID,
				UserID:      category.UserID,
				Name:        category.Name,
				Description: category.Description,
				ImageID:     category.ImageID,
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
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	for _, item := range request.Body {
		categories = append(categories, &category{
			UserID:      request.UserID,
			Name:        item.Name,
			Description: item.Description,
			ImageID:     item.ImageID,
		})
	}

	if createdCategories, err := api.interactor.createCategories(categories); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		categoriesResponse := make([]*categoryResponse, 0)

		for _, createdCategory := range createdCategories {
			categoryResponse := &categoryResponse{
				CategoryID:  createdCategory.CategoryID,
				UserID:      createdCategory.UserID,
				Name:        createdCategory.Name,
				Description: createdCategory.Description,
				ImageID:     createdCategory.ImageID,
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
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if updatedCategory, err := api.interactor.updateCategory(
		&category{
			UserID:      request.UserID,
			CategoryID:  request.CategoryID,
			Name:        request.Body.Name,
			Description: request.Body.Description,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if updatedCategory == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, categoryResponse{
			CategoryID:  updatedCategory.CategoryID,
			UserID:      updatedCategory.UserID,
			Name:        updatedCategory.Name,
			Description: updatedCategory.Description,
			ImageID:     updatedCategory.ImageID,
			CreatedAt:   updatedCategory.CreatedAt.String(),
			UpdatedAt:   updatedCategory.UpdatedAt.String(),
		})
	}
}

func (api *apiWeb) deleteCategoryHandler(ctx echo.Context) error {
	request := deleteCategoryRequest{
		UserID:     ctx.Param("user_id"),
		CategoryID: ctx.Param("category_id"),
	}

	if err := api.interactor.deleteCategory(request.UserID, request.CategoryID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

type getImagesRequest struct {
	UserID string `json:"user_id" validate:"ui"`
}

type getImageRequest struct {
	UserID  string `json:"user_id" validate:"ui"`
	ImageID string `json:"image_id" validate:"ui"`
}

type getImageRawRequest struct {
	UserID  string `json:"user_id" validate:"ui"`
	ImageID string `json:"image_id" validate:"ui"`
}

type createImagesRequest struct {
	UserID string `json:"user_id" validate:"ui"`
	Body   struct {
		Name        string `json:"name" validate:"nonzero"`
		Description string `json:"description" validate:"nonzero"`
		Url         string `json:"url"`
		ImageKey    string `json:"image_key"`
	} `json:"images" validate:"min=1"`
}

type updateImageRequest struct {
	ImageID string `json:"image_id" validate:"ui"`
	UserID  string `json:"user_id" validate:"ui"`
	Body    struct {
		Name        string `json:"name" validate:"nonzero"`
		Description string `json:"description"`
		Url         string `json:"url"`
		ImageKey    string `json:"image_key"`
	} `json:"body"`
}

type deleteImageRawRequest struct {
	UserID  string `json:"user_id" validate:"ui"`
	ImageID string `json:"image_id" validate:"ui"`
}

type imageResponse struct {
	ImageID     string `json:"image_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	FileName    string `json:"file_name,omitempty"`
	Format      string `json:"format,omitempty"`
	RawImage    []byte `json:"raw_image,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

func (api *apiWeb) registerRoutesForImages() error {
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id/images", api.getImagesHandler, api.auth)
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id/images/:image_id", api.getImageHandler, api.auth)
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id/images/:image_id/raw", api.getImageRawHandler, api.auth)
	api.client.AddRoute(http.MethodPost, "/api/1/users/:user_id/images", api.createImageHandler, api.auth)
	api.client.AddRoute(http.MethodPut, "/api/1/users/:user_id/images/:image_id", api.updateImageHandler, api.auth)
	api.client.AddRoute(http.MethodDelete, "/api/1/users/:user_id/images/:image_id", api.deleteImageHandler, api.auth)

	return nil
}

func (api *apiWeb) getImagesHandler(ctx echo.Context) error {
	request := getImagesRequest{
		UserID: ctx.Param("user_id"),
	}

	if images, err := api.interactor.getImages(request.UserID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if images == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		imagesResponse := make([]*imageResponse, 0)

		for _, image := range images {
			imageResponse := &imageResponse{
				UserID:      image.UserID,
				ImageID:     image.ImageID,
				Name:        image.Name,
				Description: image.Description,
				Url:         image.Url,
				FileName:    image.FileName,
				Format:      image.Format,
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
	request := getImageRawRequest{
		UserID:  ctx.Param("user_id"),
		ImageID: ctx.Param("image_id"),
	}

	if image, err := api.interactor.getImage(request.UserID, request.ImageID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if image == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			imageResponse{
				UserID:      image.UserID,
				ImageID:     image.ImageID,
				Name:        image.Name,
				Description: image.Description,
				Url:         image.Url,
				FileName:    image.FileName,
				Format:      image.Format,
				RawImage:    image.RawImage,
				CreatedAt:   image.CreatedAt.String(),
				UpdatedAt:   image.UpdatedAt.String(),
			})
	}
}

func (api *apiWeb) getImageRawHandler(ctx echo.Context) error {
	request := getImageRawRequest{
		UserID:  ctx.Param("user_id"),
		ImageID: ctx.Param("image_id"),
	}

	if rawImage, err := api.interactor.getImageRaw(request.UserID, request.ImageID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if rawImage == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, &imageResponse{
			ImageID:  request.ImageID,
			UserID:   request.UserID,
			RawImage: rawImage,
		})
	}
}

func (api *apiWeb) createImageHandler(ctx echo.Context) error {
	request := createImagesRequest{
		UserID: ctx.Param("user_id"),
	}

	// form values
	request.Body.Name = ctx.FormValue("name")
	request.Body.Description = ctx.FormValue("description")
	request.Body.Url = ctx.FormValue("url")
	request.Body.ImageKey = "image"

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if downloads, err := download(request.Body.ImageKey, ctx); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error uploading images").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		if len(downloads) == 0 {
			newErr := errors.FromString("there is no file in the request")
			log.Error(newErr.Error())
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
		}

		format := strings.Split(downloads[0].FileName, ".")[1]

		image := &image{
			UserID:      request.UserID,
			Name:        request.Body.Name,
			Description: request.Body.Description,
			Url:         request.Body.Url,
			FileName:    downloads[0].FileName,
			Format:      format,
			RawImage:    downloads[0].Data.Bytes(),
		}

		if createdImage, err := api.interactor.createImage(image); err != nil {
			return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
		} else {
			return ctx.JSON(http.StatusCreated, &imageResponse{
				ImageID:     createdImage.ImageID,
				UserID:      createdImage.UserID,
				Name:        createdImage.Name,
				Description: createdImage.Description,
				Url:         createdImage.Url,
				FileName:    createdImage.FileName,
				Format:      createdImage.Format,
				RawImage:    createdImage.RawImage,
				CreatedAt:   createdImage.CreatedAt.String(),
				UpdatedAt:   createdImage.UpdatedAt.String(),
			})
		}
	}
}

func (api *apiWeb) updateImageHandler(ctx echo.Context) error {
	request := updateImageRequest{
		UserID:  ctx.Param("user_id"),
		ImageID: ctx.Param("image_id"),
	}

	// form values
	request.Body.Name = ctx.FormValue("name")
	request.Body.Description = ctx.FormValue("description")
	request.Body.Url = ctx.FormValue("url")
	request.Body.ImageKey = ctx.FormValue("image_key")

	if err := validator.Validate(request); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if downloads, err := download(request.Body.ImageKey, ctx); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error uploading images").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	} else {
		if len(downloads) == 0 {
			newErr := errors.FromString("there is no file in the request")
			log.Error(newErr.Error())
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
		}

		format := strings.Split(downloads[0].FileName, ".")[1]

		image := &image{
			ImageID:     request.ImageID,
			UserID:      request.UserID,
			Name:        request.Body.Name,
			Description: request.Body.Description,
			Url:         request.Body.Url,
			Format:      format,
			FileName:    downloads[0].FileName,
			RawImage:    downloads[0].Data.Bytes(),
		}

		if updatedImage, err := api.interactor.updateImage(image); err != nil {
			return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
		} else if updatedImage == nil {
			return ctx.NoContent(http.StatusNotFound)
		} else {
			return ctx.JSON(http.StatusCreated, imageResponse{
				ImageID:     updatedImage.ImageID,
				UserID:      updatedImage.UserID,
				Name:        updatedImage.Name,
				Description: updatedImage.Description,
				Url:         updatedImage.Url,
				FileName:    updatedImage.FileName,
				Format:      updatedImage.Format,
				RawImage:    updatedImage.RawImage,
				CreatedAt:   updatedImage.CreatedAt.String(),
				UpdatedAt:   updatedImage.UpdatedAt.String(),
			})
		}
	}
}

func (api *apiWeb) deleteImageHandler(ctx echo.Context) error {
	request := deleteImageRawRequest{
		UserID:  ctx.Param("user_id"),
		ImageID: ctx.Param("image_id"),
	}

	if err := api.interactor.deleteImage(request.UserID, request.ImageID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}

type getTransactionsRequest struct {
	UserID   string `json:"user_id" validate:"ui"`
	WalletID string `json:"wallet_id" validate:"ui"`
}

type getTransactionRequest struct {
	UserID        string `json:"user_id" validate:"ui"`
	WalletID      string `json:"wallet_id" validate:"ui"`
	TransactionID string `json:"transaction_id" validate:"ui"`
}

type createTransactionsRequest struct {
	UserID   string                   `json:"user_id" validate:"ui"`
	WalletID string                   `json:"wallet_id" validate:"ui"`
	Body     []transactionItemRequest `json:"transactions" validate:"min=1"`
}

type updateTransactionRequest struct {
	UserID        string `json:"user_id" validate:"ui"`
	WalletID      string `json:"wallet_id" validate:"ui"`
	TransactionID string `json:"transaction_id" validate:"ui"`
	Body          transactionItemRequest
}

type transactionItemRequest struct {
	CategoryID  string `json:"category_id" validate:"ui"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Date        string `json:"date" validate:"nonzero"`
}

type deleteTransactionRequest struct {
	UserID        string `json:"user_id" validate:"ui"`
	WalletID      string `json:"wallet_id" validate:"ui"`
	TransactionID string `json:"transaction_id" validate:"ui"`
}

type transactionResponse struct {
	UserID        string `json:"user_id"`
	WalletID      string `json:"wallet_id"`
	TransactionID string `json:"transaction_id"`
	CategoryID    string `json:"category_id"`
	Price         string `json:"price"`
	Description   string `json:"description,omitempty"`
	Date          string `json:"date"`
	UpdatedAt     string `json:"updated_at"`
	CreatedAt     string `json:"created_at"`
}

func (api *apiWeb) registerRoutesForTransactions() error {
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id/wallets/:wallet_id/transactions", api.getTransactionsHandler, api.auth)
	api.client.AddRoute(http.MethodGet, "/api/1/users/:user_id/wallets/:wallet_id/transactions/:transaction_id", api.getTransactionHandler, api.auth)
	api.client.AddRoute(http.MethodPost, "/api/1/users/:user_id/wallets/:wallet_id/transactions", api.createTransactionsHandler, api.auth)
	api.client.AddRoute(http.MethodPut, "/api/1/users/:user_id/wallets/:wallet_id/transactions/:transaction_id", api.updateTransactionHandler, api.auth)
	api.client.AddRoute(http.MethodDelete, "/api/1/users/:user_id/wallets/:wallet_id/transactions/:transaction_id", api.deleteTransactionHandler, api.auth)

	return nil
}

func (api *apiWeb) getTransactionsHandler(ctx echo.Context) error {
	request := getTransactionsRequest{
		UserID:   ctx.Param("user_id"),
		WalletID: ctx.Param("wallet_id"),
	}

	if transactions, err := api.interactor.getTransactions(request.UserID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if transactions == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		transactionsResponse := make([]*transactionResponse, 0)

		for _, transaction := range transactions {
			transactionResponse := &transactionResponse{
				TransactionID: transaction.TransactionID,
				UserID:        transaction.UserID,
				WalletID:      transaction.WalletID,
				CategoryID:    transaction.CategoryID,
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
	request := getTransactionRequest{
		UserID:        ctx.Param("user_id"),
		WalletID:      ctx.Param("wallet_id"),
		TransactionID: ctx.Param("transaction_id"),
	}

	if transaction, err := api.interactor.getTransaction(request.UserID, request.WalletID, request.TransactionID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if transaction == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusOK,
			transactionResponse{
				TransactionID: transaction.TransactionID,
				UserID:        transaction.UserID,
				WalletID:      transaction.WalletID,
				CategoryID:    transaction.CategoryID,
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
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	for _, item := range request.Body {
		if err := validator.Validate(item); err != nil {
			newErr := errors.NewError(err)
			log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
				Error("error when validating body request").ToErrorData(newErr)
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
		}
	}

	for _, item := range request.Body {
		date, err := time.Parse(time.RFC3339, item.Date)
		if err != nil {
			newErr := errors.NewError(err)
			log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
				Error("error getting date").ToErrorData(newErr)
			return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
		}

		price, err := decimal.NewFromString(item.Price)

		transactions = append(transactions, &transaction{
			UserID:      request.UserID,
			WalletID:    request.WalletID,
			CategoryID:  item.CategoryID,
			Price:       price,
			Description: item.Description,
			Date:        date,
		})
	}

	if createdTransactions, err := api.interactor.createTransactions(transactions); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		transactionsResponse := make([]*transactionResponse, 0)

		for _, createdTransaction := range createdTransactions {
			transactionResponse := &transactionResponse{
				TransactionID: createdTransaction.TransactionID,
				UserID:        createdTransaction.UserID,
				WalletID:      createdTransaction.WalletID,
				CategoryID:    createdTransaction.CategoryID,
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
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting body").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if err := validator.Validate(request.Body); err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error when validating body request").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	price, err := decimal.NewFromString(request.Body.Price)
	if err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting price").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	date, err := time.Parse(time.RFC3339, request.Body.Date)
	if err != nil {
		newErr := errors.NewError(err)
		log.WithFields(map[string]interface{}{"error": newErr.Error(), "cause": newErr.Cause()}).
			Error("error getting date").ToErrorData(newErr)
		return ctx.JSON(http.StatusBadRequest, errorResponse{Code: http.StatusBadRequest, Message: newErr.Error(), Cause: newErr.Cause()})
	}

	if updatedTransaction, err := api.interactor.updateTransaction(
		&transaction{
			UserID:        request.UserID,
			WalletID:      request.WalletID,
			TransactionID: request.TransactionID,
			CategoryID:    request.Body.CategoryID,
			Price:         price,
			Description:   request.Body.Description,
			Date:          date,
		}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else if updatedTransaction == nil {
		return ctx.NoContent(http.StatusNotFound)
	} else {
		return ctx.JSON(http.StatusCreated, transactionResponse{
			TransactionID: updatedTransaction.TransactionID,
			UserID:        updatedTransaction.UserID,
			WalletID:      updatedTransaction.WalletID,
			CategoryID:    updatedTransaction.CategoryID,
			Price:         updatedTransaction.Price.String(),
			Description:   updatedTransaction.Description,
			Date:          updatedTransaction.Date.String(),
			CreatedAt:     updatedTransaction.CreatedAt.String(),
			UpdatedAt:     updatedTransaction.UpdatedAt.String(),
		})
	}
}

func (api *apiWeb) deleteTransactionHandler(ctx echo.Context) error {
	request := deleteTransactionRequest{
		UserID:        ctx.Param("user_id"),
		WalletID:      ctx.Param("wallet_id"),
		TransactionID: ctx.Param("transaction_id"),
	}

	if err := api.interactor.deleteTransaction(request.UserID, request.WalletID, request.TransactionID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errorResponse{Code: http.StatusInternalServerError, Message: err.Error(), Cause: err.Cause()})
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}
