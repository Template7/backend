
.PHONY: build swagger all

build:
	go build

image:
	docker-compose -f ./build/docker-compose.yaml build

swagger:
	swag init --parseDependency --parseDepth 3

run:
	./backend

all: swagger build run
