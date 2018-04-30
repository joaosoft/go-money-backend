package gomoney

import (
	gomanager "github.com/joaosoft/go-manager/app"
)

// appConfig ...
type appConfig struct {
	GoMoney MoneyConfig `json:"gomoney"`
}

// MoneyConfig ...
type MoneyConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
	Host    string             `json:"host"`
	Db      gomanager.DBConfig `json:"db"`
	Dropbox struct {
		Enabled bool `json:"enabled"`
	} `json:"dropbox"`
}
