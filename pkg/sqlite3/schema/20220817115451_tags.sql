-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags(
	ID VARCHAR (127) PRIMARY KEY UNIQUE,
	name VARCHAR (127) NOT NULL UNIQUE,
	slug VARCHAR (127) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tags;
-- +goose StatementEnd