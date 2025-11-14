FROM golang:1.24 AS builder

WORKDIR /app/payment

# копируем go.mod/go.sum inventory-сервиса
COPY payment/go.mod payment/go.sum ./
COPY shared/ /app/shared
COPY platform/ /app/platform

RUN go mod download

# копируем код сервиса
COPY payment/ ./

# копируем миграции
COPY payment/migrations app/migrations


# билдим бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd

# финальный образ
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
COPY payment/migrations ./migrations
CMD ["./main"]
