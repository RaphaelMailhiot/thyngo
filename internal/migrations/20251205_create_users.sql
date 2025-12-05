INSERT INTO users (username, email, role, created_at, updated_at)
VALUES
    ('Admin', 'admin@example.com', 'admin', now(), now()),
    ('User One', 'user1@example.com', 'user', now(), now())
    ON CONFLICT (email) DO NOTHING;