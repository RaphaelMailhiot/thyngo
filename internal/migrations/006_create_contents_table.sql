CREATE TABLE IF NOT EXISTS contents (
    id BIGSERIAL PRIMARY KEY,
    parent_table TEXT NOT NULL CHECK (parent_table IN ('posts','projects')),
    parent_id BIGINT NOT NULL,
    ord INT NOT NULL DEFAULT 0,
    type TEXT NOT NULL CHECK (type IN ('image','title','text','code','custom')),
    block_id BIGINT NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (parent_table, parent_id, ord)
);

CREATE INDEX IF NOT EXISTS idx_contents_parent ON contents (parent_table, parent_id);
CREATE INDEX IF NOT EXISTS idx_contents_block_id ON contents (block_id);
