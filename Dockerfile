FROM golang:1.21.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o event-api ./cmd/api

FROM gcr.io/distroless/debian10

COPY --from=builder /app/event-api /app/event-api
COPY --from=builder /app/docs/swagger.json /app/swagger.json
COPY --from=builder /app/.env /app/.env

WORKDIR /app

EXPOSE 8080

CMD ["/app/event-api"]