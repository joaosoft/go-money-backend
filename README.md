# go-money
[![Build Status](https://travis-ci.org/joaosoft/go-money-backend.svg?branch=master)](https://travis-ci.org/joaosoft/go-money-backend) | [![codecov](https://codecov.io/gh/joaosoft/go-money-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/go-money-backend) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/go-money-backend)](https://goreportcard.com/report/github.com/joaosoft/go-money-backend) | [![GoDoc](https://godoc.org/github.com/joaosoft/go-money-backend?status.svg)](https://godoc.org/github.com/joaosoft/go-money-backend/app)

A project that allows you to manage your day-to-day expenses.

## Support for 
* REST service
* Authentication
* Dropbox image upload and download
* Postgres database

## Start
It starts the api on port 8082 [[here](http://localhost:8082)]
```
make start 
```

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/go-money-backend/app
```

## Usage 
This examples are available in the project at [go-money-backend/bin/launcher/main.go](https://github.com/joaosoft/go-money-backend/tree/master/bin/launcher/main.go)

```go
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
```

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
