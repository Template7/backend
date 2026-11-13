.PHONY: build generate wire swagger test lint vulncheck check run image all

# tool versions (keep in sync with .github/workflows/test.yaml and build/Dockerfile)
WIRE          = go run github.com/google/wire/cmd/wire@v0.6.0
SWAG          = go run github.com/swaggo/swag/cmd/swag@v1.16.4
GOLANGCI_LINT = go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
GOVULNCHECK   = go run golang.org/x/vuln/cmd/govulncheck@v1.9.0

build: wire
	go build -ldflags="-X 'github.com/Template7/backend/internal.CommitHash=$$(git rev-parse --short HEAD)'" -o ./bin/backend ./cmd

generate: wire swagger

wire:
	$(WIRE) gen ./cmd

swagger:
	$(SWAG) init -g ./cmd/main.go -q

test:
	go test -race -count=1 ./...

lint:
	$(GOLANGCI_LINT) run ./...

vulncheck:
	$(GOVULNCHECK) ./...

# everything CI runs
check: test lint vulncheck

run:
	./bin/backend

image:
	docker-compose -f ./build/docker-compose.yaml build

all: generate build run
