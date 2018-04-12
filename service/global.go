package gomoney

import (
	logger "github.com/joaosoft/go-log/service"
)

var global = make(map[string]interface{})
var log = logger.NewLogDefault("go-money", logger.InfoLevel)
