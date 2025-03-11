-- +goose Up
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "spu_to_sku";
DROP TABLE IF EXISTS "sku_specs";
DROP TABLE IF EXISTS "sku_attr";
DROP TABLE IF EXISTS "sku";
DROP TABLE IF EXISTS "sd_product";
-- +goose StatementEnd
