package gomoney

import (
	"github.com/joaosoft/go-log/service"
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
func WithLogLevel(level golog.Level) goMoneyOption {
	return func(gomoney *goMoney) {
		log.SetLevel(level)
	}
}
