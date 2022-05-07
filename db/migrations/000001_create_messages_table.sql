-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    ts INT NOT NULL,
    text TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
