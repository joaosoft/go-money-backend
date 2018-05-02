package gomoney

import (
	"fmt"

	golog "github.com/joaosoft/go-log/app"
	gomanager "github.com/joaosoft/go-manager/app"
	"gopkg.in/validator.v2"
)

func init() {
	validator.SetValidationFunc("ui", func(v interface{}, param string) error {
		switch v.(type) {
		case string:
			if err := valUI(v.(string)); err != nil {
				return fmt.Errorf("%s is not a valid unique identifier", param)
			}
		}
		return nil
	})
}

// Money ...
type Money struct {
	interactor    *interactor
	pm            *gomanager.Manager
	config        *MoneyConfig
	isLogExternal bool
}

// NewMoney ...
func NewMoney(options ...moneyOption) (*Money, error) {
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
