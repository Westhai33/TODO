-- +goose Up
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_user_id;  -- Удаление внешнего ключа
CREATE TABLE tasks (
                       id BIGSERIAL PRIMARY KEY,              -- Уникальный идентификатор задачи
                       user_id BIGINT NOT NULL,               -- ID пользователя
                       title TEXT NOT NULL,                   -- Название задачи
                       note TEXT DEFAULT '',                  -- Заметка
                       done BOOLEAN DEFAULT FALSE,            -- Статус выполнения задачи
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Дата создания
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Дата обновления
);

-- +goose Down
DROP TABLE IF EXISTS tasks;