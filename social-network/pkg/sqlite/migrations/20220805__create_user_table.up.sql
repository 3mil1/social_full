-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS user
(
    id            VARCHAR(36)  NOT NULL PRIMARY KEY,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name    VARCHAR(255) NOT NULL,
    last_name     VARCHAR(255) NOT NULL,
    birthday      date         not null,
    image         VARCHAR(255),
    nickname      VARCHAR(255),
    about         TEXT,
    created_at    timestamp    not null default (datetime('now','localtime')),
    updated_at    timestamp    default 0,
    is_private    BOOLEAN NOT NULL CHECK (is_private IN (0, 1))
);
