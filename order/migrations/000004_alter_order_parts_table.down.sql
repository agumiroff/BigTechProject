-- +goose Up
ALTER TABLE order_parts DROP CONSTRAINT order_parts_pkey;
ALTER TABLE order_parts
    ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();

-- +goose Down
ALTER TABLE order_parts DROP COLUMN id;
ALTER TABLE order_parts
    ADD PRIMARY KEY (order_uuid, part_uuid);

