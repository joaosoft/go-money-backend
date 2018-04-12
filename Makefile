env:
	docker-compose up -d postgres manager

start: env
	docker-compose up -d app

run:
	go run ./bin/launcher/main.go

build:
	go build .

fmt:
	go fmt ./...

vet:
	go vet ./*

gometalinter:
	gometalinter ./*