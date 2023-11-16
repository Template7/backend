
.PHONY: build swagger all

build:
	go build -o ./bin/backend ./cmd

image:
	docker-compose -f ./build/docker-compose.yaml build

swagger:
	swag init -g ./cmd/main.go --parseDependency --parseDepth 3

run:
	./bin/backend

all: swagger build run
