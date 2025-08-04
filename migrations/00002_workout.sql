-- +goose up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workouts (
    id BIGSERIAL PRIMARY KEY,
    --user_id
    title VARCHAR(100) NOT NULL,
    description TEXT,
    duration_minutes INT NOT NULL, -- in minutes
    calories_burned INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose down
-- +goose StatementBegin
DROP TABLE IF EXISTS workouts;
-- +goose StatementEnd