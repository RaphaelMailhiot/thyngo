-- +goose Up
CREATE TABLE IF NOT EXISTS resource_accesses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resource_type TEXT NOT NULL CHECK (resource_type IN ('media', 'posts', 'projects', 'resumes')),
    resource_id BIGINT NOT NULL,
    permission TEXT NOT NULL DEFAULT 'view' CHECK (permission IN ('view', 'edit')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, resource_type, resource_id)
);

-- +goose Down
-- (no down migration provided)

