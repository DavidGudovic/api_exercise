-- +goose Up
-- +goose StatementBegin
ALTER TABLE workouts ADD COLUMN user_id INT NOT NULL REFERENCES public.users(id) ON DELETE CASCADE DEFAULT 1;
CREATE INDEX idx_workouts_user_id ON workouts(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workouts DROP COLUMN user_id;
DROP INDEX idx_workouts_user_id;
-- +goose StatementEnd