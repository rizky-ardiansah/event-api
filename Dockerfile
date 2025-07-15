FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o event-api ./cmd/api

FROM alpine:latest

COPY --from=builder /app/event-api /app/event-api
COPY --from=builder /app/docs/swagger.json /app/swagger.json

WORKDIR /app

EXPOSE 8080

CMD ["/app/event-api"]