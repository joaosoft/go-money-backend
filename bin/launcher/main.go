package main

import (
	"go-money/service"
	"time"

	"github.com/joaosoft/go-log/service"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "github.com/lib/pq"              // postgres driver
)

var log = golog.NewLogDefault("go-money", golog.InfoLevel)

func main() {
	start := time.Now()
	//
	// money
	app, err := gomoney.NewGoMoney()
	if err != nil {
		log.Error(err)
	} else {
		app.Start()
	}

	elapsed := time.Since(start)
	log.Infof("ELAPSED TIME: %s", elapsed)
}
