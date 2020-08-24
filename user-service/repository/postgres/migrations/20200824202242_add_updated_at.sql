-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE "user"
    ADD COLUMN "updated_at" timestamp without time zone;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE "user"
    DROP COLUMN "updated_at";