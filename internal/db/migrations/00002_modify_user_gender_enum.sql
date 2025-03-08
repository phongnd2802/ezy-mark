-- +goose Up
-- +goose StatementBegin
CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other');

ALTER TABLE "user_profile" 
    ALTER COLUMN "user_gender" TYPE gender_enum 
    USING 
        CASE 
            WHEN user_gender IS NULL THEN NULL
            WHEN user_gender = TRUE THEN 'male'::gender_enum
            ELSE 'female'::gender_enum
        END;

ALTER TABLE user_profile ALTER COLUMN user_gender SET DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user_profile" 
    ALTER COLUMN "user_gender" TYPE BOOLEAN 
    USING 
        CASE 
            WHEN user_gender = 'male' THEN TRUE
            WHEN user_gender = 'female' THEN FALSE
            ELSE NULL
        END;

DROP TYPE gender_enum;
-- +goose StatementEnd
