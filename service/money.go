package gomoney

import (
	"fmt"

	"github.com/joaosoft/go-log/service"
	"github.com/joaosoft/go-manager/service"
)

// GoMoney ...
type GoMoney struct {
	interactor *Interactor
	pm         *gomanager.GoManager
	config     *AppConfig
}

// NewGoMoney ...
func NewGoMoney(options ...GoMoneyOption) (*GoMoney, error) {
	pm := gomanager.NewManager(gomanager.WithLogger(log), gomanager.WithRunInBackground(false))

	// load configuration file
	appConfig := &AppConfig{}
	if simpleConfig, err := gomanager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := golog.ParseLevel(appConfig.Log.Level)
		log.Debugf("setting log level to %s", level)
		WithLogLevel(level)
	}
	simpleDB := gomanager.NewSimpleDB(&appConfig.Db)
	pm.AddDB("db_postgres", simpleDB)

	money := &GoMoney{
		interactor: NewInteractor(NewStorage(simpleDB), appConfig),
		pm:         pm,
		config:     appConfig,
	}

	money.Reconfigure(options...)

	return money, nil
}

// Start ...
func (api *GoMoney) Start() error {
	apiWeb := newApiWeb(api.config.Host, api.interactor)
	api.pm.AddWeb("api_web", apiWeb.new())

	return api.pm.Start()
}

// Stop ...
func (api *GoMoney) Stop() error {
	return api.pm.Stop()
}
