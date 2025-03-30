-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid (),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM USERS WHERE email = $1;

-- name: GetUserByRereshToken :one
SELECT *
FROM users
WHERE users.id = (
    SELECT user_id
    FROM refresh_tokens
    WHERE token = $1
);

-- name: UpdateUserByToken :one
UPDATE users 
SET 
    updated_at = NOW(),
    email = $2,
    hashed_password = $3
WHERE id = $1
RETURNING *;