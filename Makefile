
.PHONY: build all

build:
	go build

swagger:
	swag init

run:
	./backend

all: swagger build run
