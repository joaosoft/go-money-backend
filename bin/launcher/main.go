package main

import (
	"go-money/service"
	"time"

	"github.com/joaosoft/go-log/service"
)

var log = golog.NewLogDefault("go-money", golog.InfoLevel)

func main() {
	start := time.Now()
	//
	// money
	app := gomoney.NewGoMoney()
	app.Start()

	elapsed := time.Since(start)
	log.Infof("ELAPSED TIME: %s", elapsed)
}
