
.PHONY: build swagger all

build:
	go build

swagger:
	swag init

run:
	./backend

all: swagger build run
