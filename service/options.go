package gomoney

import "github.com/joaosoft/go-log/service"

// moneyOption ...
type moneyOption func(money *Money)

// Reconfigure ...
func (money *Money) Reconfigure(options ...moneyOption) {
	for _, option := range options {
		option(money)
	}
}

// WithLogger ...
func WithLogger(logger golog.ILog) moneyOption {
	return func(manager *Money) {
		log = logger
		manager.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level golog.Level) moneyOption {
	return func(money *Money) {
		log.SetLevel(level)
	}
}
