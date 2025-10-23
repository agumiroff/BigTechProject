FROM golang:1.24 AS builder

WORKDIR /app

# копируем go.mod/go.sum всех зависимых модулей
COPY shared/go.mod shared/go.sum ./shared/
COPY payment/go.mod payment/go.sum ./payment/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY order/go.mod order/go.sum ./order/

# грузим зависимости для order
WORKDIR /app/order
RUN go mod download

WORKDIR /app
# копируем код
COPY shared ./shared
COPY payment ./payment
COPY inventory ./inventory
COPY order ./order

WORKDIR /app/order
# билдим order
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd

# финальный образ
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
COPY order/migrations ./migrations
CMD ["./main"]
