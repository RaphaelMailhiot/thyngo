-- +goose Up
ALTER TABLE media DROP CONSTRAINT IF EXISTS media_type_check;
ALTER TABLE media ADD CONSTRAINT media_type_check CHECK (type IN ('image', 'video', 'document', 'other'));

-- +goose Down
ALTER TABLE media DROP CONSTRAINT IF EXISTS media_type_check;
ALTER TABLE media ADD CONSTRAINT media_type_check CHECK (type IN ('movie','series','other'));
