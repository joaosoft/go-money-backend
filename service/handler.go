package gomoney

import (
	"net/http"

	"github.com/joaosoft/go-manager/service"
	"github.com/labstack/echo"
)

// Handler ...
type handler struct {
	pm *gomanager.GoManager
}

// NewHandler ...
func newHandler(pm *gomanager.GoManager) *handler {
	return &handler{
		pm: pm,
	}
}

func (handler *handler) getAccountHandler(ctx echo.Context) error {
	type example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	return ctx.JSON(http.StatusOK, example{Id: ctx.Param("id"), Name: "joao", Age: 29})
}

func (handler *handler) createAccountHandler(ctx echo.Context) error {
	type example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	return ctx.JSON(http.StatusOK, example{Id: ctx.Param("id"), Name: "joao", Age: 29})
}
