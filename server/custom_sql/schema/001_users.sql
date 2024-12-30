-- +goose Up 
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY NOT NULL,
    user_name VARCHAR(20) NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    user_password TEXT NOT NULL,
    age INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;