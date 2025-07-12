-- +goose Up
-- +goose StatementBegin
-- Enable uuid-ossp extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- First, add a new UUID column to users table
ALTER TABLE users ADD COLUMN id_uuid UUID DEFAULT uuid_generate_v4();

-- Update the UUID column with generated UUIDs for existing records
UPDATE users SET id_uuid = uuid_generate_v4() WHERE id_uuid IS NULL;

-- Store the mapping between old bigint IDs and new UUIDs for sessions table
CREATE TEMP TABLE user_id_mapping AS 
SELECT id as old_id, id_uuid as new_id FROM users;

-- Add new UUID column to sessions table
ALTER TABLE sessions ADD COLUMN user_id_uuid UUID;

-- Update sessions table with new UUID user_id
UPDATE sessions SET user_id_uuid = (
    SELECT new_id FROM user_id_mapping WHERE old_id = sessions.user_id
);

-- Drop foreign key constraint on sessions
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS fk_sessions_user_id;

-- Drop the old bigint columns
ALTER TABLE sessions DROP COLUMN user_id;
ALTER TABLE users DROP COLUMN id;

-- Rename UUID columns to be the primary columns
ALTER TABLE users RENAME COLUMN id_uuid TO id;
ALTER TABLE sessions RENAME COLUMN user_id_uuid TO user_id;

-- Make the new UUID columns NOT NULL
ALTER TABLE users ALTER COLUMN id SET NOT NULL;
ALTER TABLE sessions ALTER COLUMN user_id SET NOT NULL;

-- Add primary key constraint back
ALTER TABLE users ADD PRIMARY KEY (id);

-- Add foreign key constraint back
ALTER TABLE sessions ADD CONSTRAINT fk_sessions_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add index for sessions user_id
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

-- Also need to update sessions table ID to UUID
ALTER TABLE sessions ADD COLUMN id_uuid UUID DEFAULT uuid_generate_v4();
UPDATE sessions SET id_uuid = uuid_generate_v4() WHERE id_uuid IS NULL;
ALTER TABLE sessions DROP COLUMN id;
ALTER TABLE sessions RENAME COLUMN id_uuid TO id;
ALTER TABLE sessions ALTER COLUMN id SET NOT NULL;
ALTER TABLE sessions ADD PRIMARY KEY (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- This down migration is complex and potentially destructive
-- For safety, we'll just log that this operation is not supported
SELECT 'WARNING: Down migration for UUID conversion is not implemented for data safety' as warning;
-- +goose StatementEnd
