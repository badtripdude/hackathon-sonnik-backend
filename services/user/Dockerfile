FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

FROM alpine:latest

RUN apk add --no-cache bash postgresql-client jq

WORKDIR /app

COPY --from=builder /app/app .
COPY scripts/wait_for_postgres.sh /wait_for_postgres.sh
RUN chmod +x /wait_for_postgres.sh

CMD ["./app"]
