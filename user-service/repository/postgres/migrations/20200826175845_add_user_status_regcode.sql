-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE "user"
    ADD COLUMN status character varying(30);

ALTER TABLE "user"
    ADD COLUMN registration_code character varying(6);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE "user"
    DROP COLUMN status;

ALTER TABLE "user"
    DROP COLUMN registration_code;