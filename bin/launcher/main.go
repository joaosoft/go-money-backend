package main

import (
	gomoney "github.com/joaosoft/go-money-backend/app"
	"time"

	"github.com/joaosoft/logger"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "github.com/lib/pq"              // postgres driver
)

var log = logger.NewLogDefault("go-money", logger.InfoLevel)

func main() {
	start := time.Now()
	//
	// money
	app, err := gomoney.NewMoney()
	if err != nil {
		log.Error(err)
	} else {
		app.Start()
	}

	elapsed := time.Since(start)
	log.Infof("ELAPSED TIME: %s", elapsed)
}
