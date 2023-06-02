package gomoney

import (
	"fmt"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
	"github.com/joaosoft/validator"
)

func init() {
	validator.AddCallback("ui", func(context *validator.ValidatorContext, validationData *validator.ValidationData) []error {
		switch v := validationData.Value.Interface().(type) {
		case string:
			if err := valUI(v); err != nil {
				return []error{fmt.Errorf("%s is not a valid unique identifier", v)}
			}
		}
		return nil
	})
}

// Money ...
type Money struct {
	interactor    *interactor
	pm            *manager.Manager
	config        *MoneyConfig
	isLogExternal bool
}

// NewMoney ...
func NewMoney(options ...moneyOption) (*Money, error) {
	pm := manager.NewManager(manager.WithRunInBackground(false))

	money := &Money{
		pm: pm,
	}

	money.Reconfigure(options...)

	if money.isLogExternal {
		pm.Reconfigure(manager.WithLogger(log))
	}

	// load configuration file
	appConfig := &appConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", getEnv()), appConfig); err != nil {
		log.Error(err.Error())
	} else {
		pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.GoMoney.Log.Level)
		log.Debugf("setting log level to %s", level)
		log.Reconfigure(logger.WithLevel(level))
	}

	simpleDB := pm.NewSimpleDB(&appConfig.GoMoney.Db)
	if err := pm.AddDB("db_postgres", simpleDB); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	money.config = &appConfig.GoMoney
	money.interactor = newInteractor(newStoragePostgres(simpleDB), newStorageDropbox(nil), &appConfig.GoMoney)

	return money, nil
}

// Start ...
func (m *Money) Start() error {
	apiWeb := m.newApiWeb(m.config.Host, m.interactor)
	m.pm.AddWeb("api_web", apiWeb.client)

	return m.pm.Start()
}

// Stop ...
func (m *Money) Stop() error {
	return m.pm.Stop()
}
