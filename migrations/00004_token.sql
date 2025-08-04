-- +goose up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tokens (
    hash BYTEA PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry TIMESTAMPTZ NOT NULL,
    scope VARCHAR(50) NOT NULL
);
-- +goose StatementEnd

-- +goose down
-- +goose StatementBegin
DROP TABLE IF EXISTS tokens;
-- +goose StatementEnd