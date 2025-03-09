-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "roles" (
    "role_id" serial PRIMARY KEY,
    "role_name" varchar NOT NULL UNIQUE,
    "description" text DEFAULT NULL,
    "created_at" timestamptz  NOT NULL DEFAULT (now()),
    "updated_at" timestamptz  NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "permissions" (
    "permission_id" serial PRIMARY KEY,
    "permission_name" varchar NOT NULL UNIQUE,
    "description" text DEFAULT NULL,
    "created_at" timestamptz  NOT NULL DEFAULT (now()),
    "updated_at" timestamptz  NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "user_roles" (
    "user_id" bigint NOT NULL,
    "role_id" integer NOT NULL,
    "assigned_at" timestamptz  NOT NULL DEFAULT (now()),
    PRIMARY KEY ("user_id", "role_id"),
    FOREIGN KEY ("user_id") REFERENCES "user_base" ("user_id") ON DELETE CASCADE,
    FOREIGN KEY ("role_id") REFERENCES "roles" ("role_id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "role_permissions" (
    "role_id" integer NOT NULL,
    "permission_id" integer NOT NULL,
    "granted_at" timestamptz  NOT NULL DEFAULT (now()),
    PRIMARY KEY ("role_id", "permission_id"),
    FOREIGN KEY ("role_id") REFERENCES "roles" ("role_id") ON DELETE CASCADE,
    FOREIGN KEY ("permission_id") REFERENCES "permissions" ("permission_id") ON DELETE CASCADE
);

-- Seed default data
INSERT INTO "roles" ("role_name", "description") VALUES 
    ('admin', 'Has full control over the system, including managing users, shops, orders, and system settings.'),
    ('shop', 'Manages a specific shop, including products, orders, and employees.'),
    ('customer', 'A registered user who can browse, purchase, and review products.'),
    ('employee', 'Handles customer support and order processing, assigned by a shop owner.');

-- Seed default permissions
INSERT INTO "permissions" ("permission_name", "description") VALUES
    ('manage_users', 'Create, update, delete users and assign roles.'),
    ('manage_shops', 'Create, update, delete shops and manage shop settings.'),
    ('manage_inventory', 'Manage product inventory, stock levels, and pricing.'),
    ('customer_support', 'Handle customer inquiries and complaints.'),
    ('write_reviews', 'Write and edit product reviews.'),
    ('view_products', 'Browse and view product listings.');

-- Admin permissions
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES 
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'admin'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'manage_users')),
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'admin'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'manage_shops')),
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'admin'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'customer_support')),
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'admin'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'view_products'));


-- Shop permissions
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES 
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'shop'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'manage_shops')),
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'shop'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'manage_inventory')),
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'shop'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'customer_support')),
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'shop'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'view_products'));

-- Customer permissions
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES 
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'customer'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'write_reviews')),
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'customer'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'view_products'));

-- Employee permissions
INSERT INTO "role_permissions" ("role_id", "permission_id") VALUES 
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'employee'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'customer_support')),
    ((SELECT "role_id" FROM "roles" WHERE "role_name" = 'employee'), (SELECT "permission_id" FROM "permissions" WHERE "permission_name" = 'view_products'));


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "role_permissions";
DROP TABLE IF EXISTS "user_roles";
DROP TABLE IF EXISTS "permissions";
DROP TABLE IF EXISTS "roles";
-- +goose StatementEnd
