package gomoney

import (
	"github.com/joaosoft/go-manager/service"
)

// appConfig ...
type appConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Host string             `json:"host"`
	Db   gomanager.DBConfig `json:"db"`
}
