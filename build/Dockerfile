FROM golang:alpine as builder
ARG PROJECT=backend

WORKDIR /$PROJECT
COPY . .
RUN go mod download

# install swag cli then build swagger document
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init --parseDependency --parseDepth 3

RUN go build -o $PROJECT .

FROM alpine
ENV PROJECT=backend

WORKDIR /$PROJECT
COPY --from=builder /$PROJECT/$PROJECT ./
COPY ./resource/template ./resource/template

CMD ["sh", "-c", "./$PROJECT"]
