package gomoney

import (
	"github.com/joaosoft/go-manager/service"
)

// appConfig ...
type appConfig struct {
	GoMoney goMoneyConfig `json:"gomoney"`
}

// goMoneyConfig ...
type goMoneyConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Host string             `json:"host"`
	Db   gomanager.DBConfig `json:"db"`
}
