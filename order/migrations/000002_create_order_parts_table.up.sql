-- +goose Up

  CREATE TABLE IF NOT EXISTS order_parts (
    order_uuid UUID NOT NULL REFERENCES orders(order_uuid) ON DELETE CASCADE,
    part_uuid UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    PRIMARY KEY (order_uuid, part_uuid),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_order_parts_order_uuid ON order_parts(order_uuid);
CREATE INDEX IF NOT EXISTS idx_order_parts_part_uuid ON order_parts(part_uuid);


-- +goose Down

DROP INDEX IF EXISTS idx_order_parts_order_uuid;
DROP INDEX IF EXISTS idx_order_parts_part_uuid;
DROP TABLE IF EXISTS order_parts;
