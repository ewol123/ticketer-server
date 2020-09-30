-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE "user" (
  "id" uuid PRIMARY KEY,
  "created_at" timestamp,
  "full_name" varchar,
  "email" varchar UNIQUE,
  "password" varchar
);

CREATE TABLE "role" (
  "id" uuid PRIMARY KEY,
  "name" varchar UNIQUE
);

CREATE TABLE "user_role" (
  "user_id" uuid,
  "role_id" uuid
);


ALTER TABLE "user_role" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "user_role" ADD FOREIGN KEY ("role_id") REFERENCES "role" ("id");

CREATE UNIQUE INDEX ON "user" ("id");
CREATE UNIQUE INDEX ON "role" ("id");
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE "user" CASCADE;
DROP TABLE  "role" CASCADE;
DROP TABLE  "user_role" CASCADE;