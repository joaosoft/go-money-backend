package goaccount

import (
	"github.com/joaosoft/go-manager/service"
)

// GoAccount ...
type GoAccount struct {
	pm *gomanager.GoManager
}

// NewGoAccount ...
func NewGoAccount(options ...GoAccountOption) *GoAccount {
	account := &GoAccount{
		pm: gomanager.NewManager(gomanager.WithLogger(log)),
	}
	account.Reconfigure(options...)

	return account
}

func (goaccount *GoAccount) Start() {
	handler := newHandler(goaccount.pm)
	web := gomanager.NewSimpleWebEcho("8080")
	web.AddRoute("GET", "/account/:id", handler.getAccountHandler)
	web.AddRoute("POST", "/account", handler.createAccountHandler)

	goaccount.pm.AddWeb("web_account", web)
}

func (service *GoAccount) Stop() {
	service.Stop()
}
