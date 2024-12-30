-- +goose Up 
ALTER TABLE users
ADD COLUMN user_type VARCHAR(10) DEFAULT 'user',
ADD CONSTRAINT chk_user_type CHECK (user_type IN ('user', 'admin'));

-- +goose Down
ALTER TABLE users
DROP COLUMN user_type
DROP CONSTRAINT chk_user_type 