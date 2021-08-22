FROM golang:alpine as builder
ARG PROJECT=backend

WORKDIR /$PROJECT
COPY . .
RUN go mod download
RUN go build -o $PROJECT .

FROM alpine
ARG PROJECT=backend

WORKDIR /$PROJECT
COPY --from=builder /$PROJECT/$PROJECT ./
COPY ./resource/template ./resorce/template

CMD ["./$PROJECT"]
