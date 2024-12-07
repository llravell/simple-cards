-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  login VARCHAR(40) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
