FROM golang:1.24 AS builder

WORKDIR /app/inventory

# копируем go.mod/go.sum inventory-сервиса
COPY inventory/go.mod inventory/go.sum ./
COPY shared/ /app/shared
COPY platform/ /app/platform

RUN go mod download

# копируем код сервиса
COPY inventory/ ./

# копируем миграции
COPY inventory/migrations app/migrations


# билдим бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd

# финальный образ
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
COPY inventory/migrations ./migrations
CMD ["./main"]
