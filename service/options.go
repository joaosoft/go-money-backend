package gomoney

import (
	logger "github.com/joaosoft/go-log/service"
)

// GoMoneyOption ...
type GoMoneyOption func(gomoney *GoMoney)

// Reconfigure ...
func (gomoney *GoMoney) Reconfigure(options ...GoMoneyOption) {
	for _, option := range options {
		option(gomoney)
	}
}

// WithLevel ...
func WithLevel(level logger.Level) GoMoneyOption {
	return func(gomoney *GoMoney) {
		log.SetLevel(level)
	}
}
