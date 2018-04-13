package gomoney

import (
	"net/http"

	"github.com/joaosoft/go-manager/service"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
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

	web.AddRoute("GET", "/user/:id", api.getUserHandler)
	web.AddRoute("POST", "/users", api.createUserHandler)
	web.AddRoute("PUT", "/user/:id", api.updateUserHandler)
	web.AddRoute("DELETE", "/user/:id", api.deleteUserHandler)

	return web
}

func (api *apiWeb) getUserHandler(ctx echo.Context) error {
	userID, err := uuid.FromString(ctx.Param("id"))
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
	userID, err := uuid.FromString(ctx.Param("id"))
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
	userID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := api.interactor.DeleteUser(userID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}
