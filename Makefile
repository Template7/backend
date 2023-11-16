
.PHONY: build swagger all

build:
	go build -o ./bin/backend ./cmd

image:
	docker-compose -f ./build/docker-compose.yaml build

swagger:
	swag init --parseDependency --parseDepth 3 -g ./cmd/main.go

run:
	./bin/backend

all: swagger build run
