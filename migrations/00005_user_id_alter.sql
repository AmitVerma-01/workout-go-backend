-- +goose up
-- +goose StatementBegin
ALTER TABLE workouts
ADD COLUMN user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose down
-- +goose StatementBegin
ALTER TABLE workouts
DROP COLUMN user_id;
-- +goose StatementEnd