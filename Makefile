env:
	docker-compose up -d postgres manager

start: env
	docker-compose up -d app

run:
	go run ./bin/launcher/main.go

utest:
	npm run utest

test:
	npm test

fmt:
	go fmt ./...

vet:
	go vet ./*

metalinter:
	gometalinter ./*

swagger:
	mkdir -p assets/swagger assets/redoc
	swagger generate spec -o assets/swagger/swagger.json -b ./app
	swagger generate spec -o assets/redoc/swagger.json -b ./app

build:
	docker build -t go-money-backend-image .

push:
	docker login
	docker tag go-money-backend-image joaosoft/go-money-backend-image
	docker push joaosoft/go-money-backend-image

dockerhub: build push