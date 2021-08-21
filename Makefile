
.PHONY: build swagger all

build:
	go build

swagger:
	swag init --parseDependency --parseDepth 3

run:
	./backend

all: swagger build run
