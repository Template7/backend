# keep the Go version in sync with .github/workflows/test.yaml
FROM golang:1.27.4-alpine AS builder
ARG PROJECT=backend
ARG COMMIT_HASH

RUN echo "commit hash: $COMMIT_HASH"

WORKDIR /$PROJECT
COPY . .
RUN go mod download

# regenerate swagger document and dependency injection (pinned tool versions)
RUN go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g ./cmd/main.go -q
RUN go run github.com/google/wire/cmd/wire@v0.7.0 gen ./cmd

RUN go build -ldflags="-X 'github.com/Template7/backend/internal.CommitHash=$COMMIT_HASH'" -o $PROJECT ./cmd

FROM alpine

ENV PROJECT=backend

WORKDIR /$PROJECT
COPY --from=builder /$PROJECT/$PROJECT ./

CMD ["sh", "-c", "./$PROJECT"]
