include .env
export $(shell sed 's/=.*//' .env)
.PHONY:
.SILENT:

build:
	go mod download && go build -o ./.bin/app ./app/cmd/main.go

run: build
	./.bin/app

publish:
	go run app/cmd/publish/publish.go

migrate-up:
	migrate -path ./migrations -database 'postgres://pguser:${DB_PASSWORD}@localhost:5431/devdb?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://pguser:${DB_PASSWORD}@localhost:5431/devdb?sslmode=disable' down

wrk-server:
	go-wrk -c 120 -d 10 http://localhost:8000/ping

wrk-orders:
	go-wrk -c 120 -d 10 http://localhost:8000/api/order

wrk: wrk-server wrk-orders