FROM golang:1.23.0-alpine

RUN apk add --no-cache bash git postgresql-client

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY entrypoint.sh /goose/entrypoint.sh
RUN chmod +x /goose/entrypoint.sh
