package gomoney

import (
	logger "github.com/joaosoft/go-log/service"
)

// goMoneyOption ...
type goMoneyOption func(gomoney *goMoney)

// reconfigure ...
func (gomoney *goMoney) reconfigure(options ...goMoneyOption) {
	for _, option := range options {
		option(gomoney)
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) goMoneyOption {
	return func(gomoney *goMoney) {
		log.SetLevel(level)
	}
}
