-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION assign_default_role()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_roles (user_id, role_id)
    VALUES (
        NEW.user_id,
        (SELECT role_id FROM roles WHERE role_name = 'customer')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_assign_default_role
AFTER INSERT ON user_base
FOR EACH ROW
EXECUTE FUNCTION assign_default_role();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_assign_default_role ON user_base;
DROP FUNCTION IF EXISTS assign_default_role;
-- +goose StatementEnd
