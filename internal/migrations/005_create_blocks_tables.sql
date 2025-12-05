CREATE TABLE IF NOT EXISTS blocks (
    id BIGSERIAL PRIMARY KEY,
    type TEXT NOT NULL CHECK (type IN ('image','title','text','code','custom')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS image_blocks (
    id BIGSERIAL PRIMARY KEY,
    block_id BIGINT NOT NULL UNIQUE REFERENCES blocks(id) ON DELETE CASCADE,
    url TEXT,
    title TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS title_blocks (
    id BIGSERIAL PRIMARY KEY,
    block_id BIGINT NOT NULL UNIQUE REFERENCES blocks(id) ON DELETE CASCADE,
    level SMALLINT,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS text_blocks (
    id BIGSERIAL PRIMARY KEY,
    block_id BIGINT NOT NULL UNIQUE REFERENCES blocks(id) ON DELETE CASCADE,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS code_blocks (
    id BIGSERIAL PRIMARY KEY,
    block_id BIGINT NOT NULL UNIQUE REFERENCES blocks(id) ON DELETE CASCADE,
    language TEXT,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS custom_blocks (
    id BIGSERIAL PRIMARY KEY,
    block_id BIGINT NOT NULL UNIQUE REFERENCES blocks(id) ON DELETE CASCADE,
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_custom_blocks_payload_gin ON custom_blocks USING gin (payload);
