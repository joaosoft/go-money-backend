package gomoney

import (
	"github.com/joaosoft/go-manager/service"
)

// GoMoney ...
type GoMoney struct {
	pm *gomanager.GoManager
}

// NewGoMoney ...
func NewGoMoney(options ...GoMoneyOption) *GoMoney {
	account := &GoMoney{
		pm: gomanager.NewManager(gomanager.WithLogger(log)),
	}
	account.Reconfigure(options...)

	return account
}

func (gomoney *GoMoney) Start() {
	handler := newHandler(gomoney.pm)
	web := gomanager.NewSimpleWebEcho("8080")
	web.AddRoute("GET", "/user/:id", handler.getAccountHandler)
	web.AddRoute("POST", "/expenses", handler.createAccountHandler)
	web.AddRoute("GET", "/expenses/:id", handler.getAccountHandler)
	web.AddRoute("POST", "/expenses", handler.createAccountHandler)

	gomoney.pm.AddWeb("web_account", web)
}

func (service *GoMoney) Stop() {
	service.Stop()
}
