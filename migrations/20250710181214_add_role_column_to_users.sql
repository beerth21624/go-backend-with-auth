-- +goose Up
-- +goose StatementBegin
-- Add role column to users table
ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';

-- Create index for role column
CREATE INDEX idx_users_role ON users(role);

-- Update existing admin user to have admin role
UPDATE users SET role = 'admin' WHERE username = 'admin';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove role column and index
DROP INDEX IF EXISTS idx_users_role;
ALTER TABLE users DROP COLUMN IF EXISTS role;
-- +goose StatementEnd
