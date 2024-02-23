-- +goose Up
-- +goose StatementBegin
CREATE TABLE images (
    id VARCHAR(36) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE images;
-- +goose StatementEnd