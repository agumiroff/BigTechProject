-- +goose Up

CREATE TABLE IF NOT EXISTS orders (
    order_uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    status TEXT NOT NULL,
    payment_method TEXT,
    transaction_uuid UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);


-- Indexes
CREATE INDEX IF NOT EXISTS idx_orders_user_uuid ON orders(user_uuid);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
-- +goose Down

DROP INDEX IF EXISTS idx_order_parts_order_uuid;
DROP INDEX IF EXISTS idx_order_parts_part_uuid;
DROP INDEX IF EXISTS idx_orders_user_uuid;
DROP INDEX IF EXISTS idx_orders_status;
DROP TABLE IF EXISTS orders;
