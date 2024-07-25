-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    user_id SERIAL PRIMARY KEY,
    passport_number VARCHAR(255) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    patronymic VARCHAR(50) NOT NULL,
    address VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks
(
    task_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id),
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    stop_time TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tasks