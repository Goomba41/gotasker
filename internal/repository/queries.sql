-- name: CreateUser :one
INSERT INTO users (email, "password")
VALUES($1, $2)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET email=$2, "password"=$3
WHERE id=$1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1
RETURNING *;
