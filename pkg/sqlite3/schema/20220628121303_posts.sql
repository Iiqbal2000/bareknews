-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS news(
	ID CHAR (127) PRIMARY KEY UNIQUE,
	title VARCHAR (127) NOT NULL UNIQUE,
	slug VARCHAR (127) NOT NULL UNIQUE,
	status VARCHAR (127) NOT NULL,
	body TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE news;
-- +goose StatementEnd
