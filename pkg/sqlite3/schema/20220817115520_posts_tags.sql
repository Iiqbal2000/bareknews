-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS news_tags(
	newsID VARCHAR (127) NOT NULL,
	tagsID VARCHAR (127) NOT NULL,
	FOREIGN KEY(newsID) REFERENCES news(id) ON DELETE CASCADE,
	FOREIGN KEY(tagsID) REFERENCES tags(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE news_tags;
-- +goose StatementEnd