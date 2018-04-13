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
	// load configuration file
	configApp := &AppConfig{}
	if _, err := readFile(fmt.Sprintf("/config/app.%s.json", getEnv()), configApp); err != nil {
		log.Error(err)
	} else {
		level, _ := golog.ParseLevel(configApp.Log.Level)
		log.Debugf("setting log level to %s", level)
		WithLogLevel(level)
	}

	conn, err := configApp.Db.Connect()
	if err != nil {
		return nil, err
	}
	money := &GoMoney{
		interactor: NewInteractor(NewStorage(conn)),
		pm:         gomanager.NewManager(gomanager.WithLogger(log), gomanager.WithRunInBackground(false)),
		config:     configApp,
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
