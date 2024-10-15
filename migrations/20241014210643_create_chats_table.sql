-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat (
    id SERIAL PRIMARY KEY,
    users TEXT NOT NULL UNIQUE    
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd
