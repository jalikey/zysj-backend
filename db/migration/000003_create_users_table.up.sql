CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(255) UNIQUE NOT NULL,
  "password_hash" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Add an index on username for faster lookups
CREATE INDEX ON "users" ("username");