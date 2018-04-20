package gomoney

import (
	"fmt"

	"github.com/joaosoft/go-log/service"
	"github.com/joaosoft/go-manager/service"
)

// goMoney ...
type goMoney struct {
	interactor *interactor
	pm         *gomanager.GoManager
	config     *appConfig
}

// NewGoMoney ...
func NewGoMoney(options ...goMoneyOption) (*goMoney, error) {
	pm := gomanager.NewManager(gomanager.WithLogger(log), gomanager.WithRunInBackground(false))

	// load configuration file
	appConfig := &appConfig{}
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

	money := &goMoney{
		interactor: newInteractor(newStoragePostgres(simpleDB), appConfig),
		pm:         pm,
		config:     appConfig,
	}

	money.reconfigure(options...)

	return money, nil
}

// Start ...
func (api *goMoney) Start() error {
	apiWeb := newApiWeb(api.config.Host, api.interactor)
	api.pm.AddWeb("api_web", apiWeb.new())

	return api.pm.Start()
}

// Stop ...
func (api *goMoney) Stop() error {
	return api.pm.Stop()
}
