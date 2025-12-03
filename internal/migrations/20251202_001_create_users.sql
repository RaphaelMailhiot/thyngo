CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO users (username, email, role, created_at, updated_at)
VALUES
    ('Admin', 'admin@example.com', 'admin', now(), now()),
    ('User One', 'user1@example.com', 'user', now(), now())
    ON CONFLICT (email) DO NOTHING;