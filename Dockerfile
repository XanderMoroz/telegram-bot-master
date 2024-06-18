#---Build stage---
FROM golang:1.13 AS builder
COPY . /go/src/telegram-bot
WORKDIR /go/src/telegram-bot/cmd/telegram-bot-server

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s' -o /go/bin/service

#---Final stage---
FROM alpine:latest
COPY --from=builder /go/bin/service /go/bin/service
CMD /go/bin/service --port 8090 --host '0.0.0.0'