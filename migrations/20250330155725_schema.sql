-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    uuid       UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(100) NOT NULL,
    created_at TIMESTAMP(3) NOT NULL,
    updated_at TIMESTAMP(3)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
