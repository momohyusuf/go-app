-- name: CreateUser :one
INSERT INTO users (
user_name, email, user_password, age, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)

RETURNING  user_name, email, created_at;

-- name: FindUserByEmail :one
SELECT user_id, user_name, email, user_password, user_type
FROM users
WHERE email = $1;

