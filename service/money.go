package gomoney

import (
	"fmt"

	"github.com/joaosoft/go-log/service"
	"github.com/joaosoft/go-manager/service"
)

// Money ...
type Money struct {
	interactor    *interactor
	pm            *gomanager.Manager
	config        *MoneyConfig
	isLogExternal bool
}

// NewGoMoney ...
func NewGoMoney(options ...moneyOption) (*Money, error) {
	pm := gomanager.NewManager(gomanager.WithRunInBackground(false))

	money := &Money{
		pm: pm,
	}

	money.Reconfigure(options...)

	if money.isLogExternal {
		pm.Reconfigure(gomanager.WithLogger(log))
	}

	// load configuration file
	appConfig := &appConfig{}
	if simpleConfig, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := golog.ParseLevel(appConfig.GoMoney.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(golog.WithLevel(level))
	}

	simpleDB := gomanager.NewSimpleDB(&appConfig.GoMoney.Db)
	if err := pm.AddDB("db_postgres", simpleDB); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	money.config = &appConfig.GoMoney
	money.interactor = newInteractor(newStoragePostgres(simpleDB), newStorageDropbox(nil), &appConfig.GoMoney)

	return money, nil
}

// Start ...
func (api *Money) Start() error {
	apiWeb := newApiWeb(api.config.Host, api.interactor)
	api.pm.AddWeb("api_web", apiWeb.client)

	return api.pm.Start()
}

// Stop ...
func (api *Money) Stop() error {
	return api.pm.Stop()
}
