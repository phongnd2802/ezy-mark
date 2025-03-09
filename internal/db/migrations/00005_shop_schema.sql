-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "shops" (
    "shop_id" bigserial PRIMARY KEY,
    "owner_id" bigint NOT NULL,
    "shop_name" varchar NOT NULL,
    "shop_description" text DEFAULT NULL,
    "shop_logo" varchar NOT NULL,
    "shop_phone" varchar DEFAULT NULL,
    "shop_email" varchar NOT NULL UNIQUE,
    "shop_address" varchar NOT NULL,
    "business_license" varchar NOT NULL,
    "tax_id" varchar NOT NULL UNIQUE,    
    "is_active" boolean DEFAULT false,
    "is_blocked" boolean DEFAULT false,
    "is_verified" boolean DEFAULT false,
    "verified_by" bigint DEFAULT NULL,
    "verified_at" timestamptz DEFAULT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    FOREIGN KEY ("owner_id") REFERENCES "user_base" ("user_id") ON DELETE CASCADE,
    FOREIGN KEY ("verified_by") REFERENCES "user_base" ("user_id") ON DELETE SET NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "shops";
-- +goose StatementEnd
