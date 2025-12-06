-- +goose Up
-- +goose StatementBegin
ALTER TABLE workout_entries ALTER COLUMN reps DROP NOT NULL;
ALTER TABLE workout_entries ALTER COLUMN sets DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workout_entries ALTER COLUMN reps SET NOT NULL;
ALTER TABLE workout_entries ALTER COLUMN sets SET NOT NULL;
-- +goose StatementEnd