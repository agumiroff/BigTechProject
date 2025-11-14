# 1️⃣ Builder stage
FROM golang:1.24 AS builder

WORKDIR /app

# копируем go.mod/go.sum всех зависимых модулей (для кэширования зависимостей)
COPY shared/go.mod shared/go.sum ./shared/
COPY payment/go.mod payment/go.sum ./payment/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY order/go.mod order/go.sum ./order/
COPY platform/go.mod platform/go.sum ./platform/

# загружаем зависимости для order
WORKDIR /app/order
RUN go mod download

# копируем весь код
WORKDIR /app
COPY shared ./shared
COPY payment ./payment
COPY inventory ./inventory
COPY order ./order
COPY platform ./platform

# сборка order
WORKDIR /app/order
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd

# 2️⃣ Final stage
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
COPY order/migrations ./migrations
CMD ["./main"]
