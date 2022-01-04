FROM golang:1.17.5-alpine AS builder

WORKDIR /github.com/ulstu-schedule/bot-vk/

COPY . .

RUN go mod download
RUN go build -o ./.bin/bot ./cmd/bot/main.go

FROM alpine:latest

# Install base packages
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates

# Change TimeZone
RUN apk add --update tzdata
ENV TZ=Europe/Samara

# Clean APK cache
RUN rm -rf /var/cache/apk/*

WORKDIR /root/

COPY --from=0 /github.com/ulstu-schedule/bot-vk/.bin/bot .
COPY --from=0 /github.com/ulstu-schedule/bot-vk/configs configs/

CMD ["./bot"]