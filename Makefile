
.PHONY: build swagger all

build:
	 wire gen ./cmd; go build -ldflags="-X 'github.com/Template7/backend/internal.CommitHash=$$(git rev-parse --short  HEAD)'" -o ./bin/backend ./cmd

image:
	docker-compose -f ./build/docker-compose.yaml build

swagger:
	swag init -g ./cmd/main.go

run:
	./bin/backend

all: swagger build run
