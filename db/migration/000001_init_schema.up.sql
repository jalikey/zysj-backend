-- 000001_init_schema.up.sql

-- Enable pgcrypto extension for UUID generation if needed later
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create the categories table
CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "slug" varchar(255) UNIQUE NOT NULL,
  "description" text,
  "parent_id" bigint REFERENCES "categories"("id") ON DELETE SET NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create the articles table
CREATE TABLE "articles" (
  "id" bigserial PRIMARY KEY,
  "title" varchar(255) NOT NULL,
  "content" text NOT NULL,
  "category_id" bigint REFERENCES "categories"("id") ON DELETE CASCADE,
  "author" varchar(255),
  "source" varchar(255),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Add indexes for performance on foreign keys and frequently queried columns
CREATE INDEX ON "categories" ("slug");
CREATE INDEX ON "categories" ("parent_id");
CREATE INDEX ON "articles" ("category_id");

-- Add a trigger to update the updated_at timestamp on articles table
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_articles_updated_at
BEFORE UPDATE ON articles
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
