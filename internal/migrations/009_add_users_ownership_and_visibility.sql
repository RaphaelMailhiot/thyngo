-- +goose Up
-- Add user_id and visibility to posts
ALTER TABLE posts ADD COLUMN IF NOT EXISTS user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS visibility TEXT NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'private', 'shared'));

-- Add user_id and visibility to projects
ALTER TABLE projects ADD COLUMN IF NOT EXISTS user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS visibility TEXT NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'private', 'shared'));

-- Add user_id and visibility to media
ALTER TABLE media ADD COLUMN IF NOT EXISTS user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE media ADD COLUMN IF NOT EXISTS visibility TEXT NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'private', 'shared'));

-- Add visibility to resumes
ALTER TABLE resumes ADD COLUMN IF NOT EXISTS visibility TEXT NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'private', 'shared'));

-- +goose Down
-- (no down migration provided)

