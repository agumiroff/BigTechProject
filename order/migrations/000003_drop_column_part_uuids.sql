-- +goose Up

ALTER TABLE orders DROP COLUMN IF EXISTS part_uuids;
