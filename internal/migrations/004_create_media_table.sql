-- +goose Up
CREATE TABLE IF NOT EXISTS media (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    slug TEXT UNIQUE,
    title TEXT,
    type TEXT NOT NULL DEFAULT 'other' CHECK (type IN ('movie','series','other')),
    visibility TEXT NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'private', 'shared')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_media_slug ON media (slug);

-- +goose Down
-- (no down migration provided)

