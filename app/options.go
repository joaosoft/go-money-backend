package gomoney

import "github.com/joaosoft/logger"

// moneyOption ...
type moneyOption func(money *Money)

// Reconfigure ...
func (m *Money) Reconfigure(options ...moneyOption) {
	for _, option := range options {
		option(m)
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) moneyOption {
	return func(manager *Money) {
		log = logger
		manager.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) moneyOption {
	return func(money *Money) {
		log.SetLevel(level)
	}
}
