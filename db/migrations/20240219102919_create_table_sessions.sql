-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id VARCHAR(36) PRIMARY KEY,
    data text NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd