-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION postgis;
CREATE TYPE "FaultType" AS ENUM ('leak');
CREATE TYPE "StatusType" AS ENUM ('draft', 'done', 'inactive');

CREATE TABLE "ticket" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid,
  "worker_id" uuid,
  "fault_type" "FaultType",
  "address" varchar,
  "full_name" varchar,
  "phone" varchar,
  "geo_location" geography,
  "image_url" varchar,
  "status" "StatusType",
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE UNIQUE INDEX ON "ticket" ("id");
CREATE INDEX ON "ticket" USING gist("geo_location");
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE "ticket" CASCADE;
DROP TYPE "FaultType" CASCADE;
DROP TYPE "StatusType" CASCADE;