FROM golang:alpine as builder
ARG PROJECT=backend
ARG COMMIT_HASH

RUN echo "commit hash: $COMMIT_HASH"

WORKDIR /$PROJECT
COPY . .
RUN go mod download

# install swag cli then build swagger document
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/google/wire/cmd/wire@latest
RUN swag init -g ./cmd/main.go
RUN wire gen ./cmd

RUN go build -ldflags="-X 'github.com/Template7/backend/internal.CommitHash=$COMMIT_HASH'" -o $PROJECT ./cmd

FROM alpine

ARG PROJECT=backend

WORKDIR /$PROJECT
COPY --from=builder /$PROJECT/$PROJECT ./

CMD ["sh", "-c", "./$PROJECT"]
