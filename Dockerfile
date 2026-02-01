FROM golang:1.25-alpine AS builder
RUN apk add --no-cache ca-certificates tzdata curl inotify-tools

WORKDIR /app

EXPOSE 8080

CMD ["sh", "-c", "/app/watch-build.sh"]