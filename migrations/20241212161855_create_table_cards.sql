-- +goose Up
-- +goose StatementBegin
CREATE TABLE cards (
  uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  term TEXT NOT NULL,
  meaning TEXT NOT NULL,
  module_uuid UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_module FOREIGN KEY(module_uuid) REFERENCES modules(uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cards;
-- +goose StatementEnd
