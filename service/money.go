package gomoney

import (
	"github.com/joaosoft/go-log/service"
	"github.com/joaosoft/go-manager/service"
)

// GoMoney ...
type GoMoney struct {
	pm     *gomanager.GoManager
	config *AppConfig
}

// NewGoMoney ...
func NewGoMoney(options ...GoMoneyOption) *GoMoney {
	// load configuration file
	configApp := &AppConfig{}
	if _, err := readFile("./config/app.json", configApp); err != nil {
		log.Error(err)
	} else {
		level, _ := golog.ParseLevel(configApp.Log.Level)
		log.Debugf("setting log level to %s", level)
		WithLogLevel(level)
	}

	account := &GoMoney{
		pm:     gomanager.NewManager(gomanager.WithLogger(log), gomanager.WithRunInBackground(false)),
		config: configApp,
	}

	account.Reconfigure(options...)

	return account
}

func (api *GoMoney) Start() error {
	webApi := newApiWeb(api.pm, api.config.Host)
	webApi.Init()

	return api.pm.Start()
}

func (api *GoMoney) Stop() error {
	return api.pm.Stop()
}
