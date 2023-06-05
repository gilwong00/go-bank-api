CREATE TABLE "users" (
  username VARCHAR PRIMARY KEY,
  hashed_password VARCHAR NOT NULL,
  first_name VARCHAR NOT NULL,
	last_name VARCHAR NOT NULL,
  email VARCHAR UNIQUE NOT NULL,
  password_changed_at TIMESTAMPTZ NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
  created_at TIMESTAMPTZ NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
