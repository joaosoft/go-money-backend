package gomoney

import (
	golog "github.com/joaosoft/logger"
)

var global = make(map[string]interface{})
var log = golog.NewLogDefault("go-money", golog.InfoLevel)

func init() {
	global[path_key] = defaultPath
}
