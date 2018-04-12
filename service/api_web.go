package gomoney

import (
	"net/http"

	"github.com/joaosoft/go-manager/service"
	"github.com/labstack/echo"
)

// apiWeb ...
type apiWeb struct {
	host string
}

// newApiWeb ...
func newApiWeb(host string) *apiWeb {
	webApi := &apiWeb{
		host: host,
	}

	return webApi
}

func (api *apiWeb) init() gomanager.IWeb {
	web := gomanager.NewSimpleWebEcho(api.host)

	web.AddRoute("GET", "/user/:id", api.getAccountHandler)
	web.AddRoute("POST", "/money", api.createAccountHandler)

	web.AddRoute("GET", "/money/:id", api.getAccountHandler)
	web.AddRoute("POST", "/money", api.createAccountHandler)

	return web
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
