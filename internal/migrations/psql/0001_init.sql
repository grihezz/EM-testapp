-- +goose Up
CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    passport_number VARCHAR(255) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    patronymic VARCHAR(50) NOT NULL,
    address VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;