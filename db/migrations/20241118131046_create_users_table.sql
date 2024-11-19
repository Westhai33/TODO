-- +goose Up
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,       -- Уникальный идентификатор пользователя
                       username TEXT NOT NULL,         -- Имя пользователя
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Дата регистрации
);

-- +goose Down
DROP TABLE IF EXISTS users;
