-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "user_base" (
    "user_id" bigserial PRIMARY KEY,
    "user_email" varchar NOT NULL,
    "user_hash" varchar NOT NULL,
    "user_password" varchar NOT NULL,
    "user_otp" varchar NOT NULL,
    "is_verified" boolean DEFAULT false,
    "is_deleted" boolean DEFAULT false,
    "created_at" timestamptz  NOT NULL DEFAULT (now()),
    "updated_at" timestamptz  NOT NULL DEFAULT (now()), 
    UNIQUE ("user_email")
);

CREATE TABLE IF NOT EXISTS "user_profile" (
    "user_id" bigint PRIMARY KEY,
    "user_email" varchar NOT NULL,
    "user_nickname" varchar NOT NULL,
    "user_fullname" varchar DEFAULT NULL,
    "user_avatar" varchar DEFAULT NULL,
    "user_mobile" varchar DEFAULT NULL,
    "user_gender" boolean DEFAULT NULL,
    "user_birthday" date DEFAULT NULL,
    "created_at" timestamptz  NOT NULL DEFAULT (now()),
    "updated_at" timestamptz  NOT NULL DEFAULT (now()),
    UNIQUE ("user_email")
);

CREATE TABLE IF NOT EXISTS "user_session" (
    "session_id" uuid PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "refresh_token" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "user_login_time" timestamptz DEFAULT NULL,
    "user_logout_time" timestamptz  DEFAULT NULL,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "user_profile" ADD FOREIGN KEY ("user_id") REFERENCES "user_base" ("user_id");
ALTER TABLE "user_session" ADD FOREIGN KEY ("user_id") REFERENCES "user_base" ("user_id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user_session";
DROP TABLE IF EXISTS "user_profile";
DROP TABLE IF EXISTS "user_base";
-- +goose StatementEnd
