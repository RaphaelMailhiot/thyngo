CREATE TABLE IF NOT EXISTS media (
    id BIGSERIAL PRIMARY KEY,
    slug TEXT UNIQUE,
    title TEXT,
    type TEXT NOT NULL DEFAULT 'other' CHECK (type IN ('movie','series','other')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_media_slug ON media (slug);
