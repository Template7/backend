FROM golang:alpine as builder
ARG PROJECT=backend

WORKDIR /$PROJECT
COPY . .
RUN go mod download

# install swag cli then build swagger document
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/google/wire/cmd/wire@latest
RUN swag init -g ./cmd/main.go
RUN wire gen ./cmd

RUN go build -o $PROJECT ./cmd

FROM alpine
ENV PROJECT=backend

WORKDIR /$PROJECT
COPY --from=builder /$PROJECT/$PROJECT ./

CMD ["sh", "-c", "./$PROJECT"]
