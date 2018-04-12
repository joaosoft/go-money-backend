package gomoney

import (
	"net/http"

	"github.com/joaosoft/go-manager/service"
	"github.com/labstack/echo"
)

// apiWeb ...
type apiWeb struct {
	host string
	pm   *gomanager.GoManager
}

// newApiWeb ...
func newApiWeb(pm *gomanager.GoManager, host string) *apiWeb {
	return &apiWeb{
		host: host,
		pm:   pm,
	}
}

func (api *apiWeb) Init() error {
	web := gomanager.NewSimpleWebEcho(api.host)
	web.AddRoute("GET", "/user/:id", api.getAccountHandler)
	web.AddRoute("POST", "/expenses", api.createAccountHandler)
	web.AddRoute("GET", "/expenses/:id", api.getAccountHandler)
	web.AddRoute("POST", "/expenses", api.createAccountHandler)

	api.pm.AddWeb("web_account", web)

	return nil
}

func (api *apiWeb) getAccountHandler(ctx echo.Context) error {
	type example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	return ctx.JSON(http.StatusOK, example{Id: ctx.Param("id"), Name: "joao", Age: 29})
}

func (api *apiWeb) createAccountHandler(ctx echo.Context) error {
	type example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	return ctx.JSON(http.StatusOK, example{Id: ctx.Param("id"), Name: "joao", Age: 29})
}
