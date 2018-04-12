package gomoney

import (
	"fmt"

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
	if _, err := readFile(fmt.Sprintf("./config/app.%s.json", getEnv()), configApp); err != nil {
		log.Error(err)
	} else {
		level, _ := golog.ParseLevel(configApp.Log.Level)
		log.Debugf("setting log level to %s", level)
		WithLogLevel(level)
	}

	money := &GoMoney{
		pm:     gomanager.NewManager(gomanager.WithLogger(log), gomanager.WithRunInBackground(false)),
		config: configApp,
	}

	money.Reconfigure(options...)

	return money
}

func (api *GoMoney) Start() error {
	apiWeb := newApiWeb(api.config.Host)
	api.pm.AddWeb("api_web", apiWeb.init())

	return api.pm.Start()
}

func (api *GoMoney) Stop() error {
	return api.pm.Stop()
}
