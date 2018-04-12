package goaccount

import (
	logger "github.com/joaosoft/go-log/service"
)

// GoAccountOption ...
type GoAccountOption func(goaccount *GoAccount)

// Reconfigure ...
func (goaccount *GoAccount) Reconfigure(options ...GoAccountOption) {
	for _, option := range options {
		option(goaccount)
	}
}

// WithLevel ...
func WithLevel(level logger.Level) GoAccountOption {
	return func(goaccount *GoAccount) {
		log.SetLevel(level)
	}
}
