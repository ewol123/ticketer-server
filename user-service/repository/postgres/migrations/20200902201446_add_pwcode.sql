-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE "user"
    ADD COLUMN reset_password_code character varying(6);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE "user"
    DROP COLUMN reset_password_code;