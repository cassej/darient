FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./app .

RUN go mod init api && \
    go mod tidy && \
    go mod download && \
    go mod verify && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /api ./main.go

#FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl

#COPY --from=builder /api /api

EXPOSE 8080

CMD ["/api"]