-- Таблица пользователей
CREATE TABLE "user" (
    login VARCHAR(16) PRIMARY KEY,
    password_hash CHAR(64) NOT NULL,  -- Строго 64 символа (SHA-256 в HEX)
    balance BIGINT NOT NULL
);

-- Таблица товаров
CREATE TABLE item (
    name VARCHAR(16) PRIMARY KEY,
    price BIGINT NOT NULL
);

-- Таблица покупок пользователя
CREATE TABLE user_item (
    id UUID PRIMARY KEY,
    user_id VARCHAR(16) NOT NULL REFERENCES "user"(login),
    item_name VARCHAR(16) NOT NULL REFERENCES item(name),
    quantity INT NOT NULL  -- Количество предметов
);

-- История переводов
CREATE TABLE coin_transfer (
    id UUID PRIMARY KEY,
    sender_id VARCHAR(16) REFERENCES "user"(login),  -- NULL, если пополнение
    receiver_id VARCHAR(16) NOT NULL REFERENCES "user"(login),
    amount BIGINT NOT NULL  -- Сумма монет
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_user_item_user_id ON user_item(user_id);
CREATE INDEX idx_user_item_item_id ON user_item(item_name);
CREATE INDEX idx_coin_transfer_sender_id ON coin_transfer(sender_id);
CREATE INDEX idx_coin_transfer_receiver_id ON coin_transfer(receiver_id);
