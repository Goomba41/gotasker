-- name: CreateUser :one
INSERT INTO users (email, "password")
VALUES($1, $2)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1
RETURNING *;

-- name: DeleteTasks :one
DELETE FROM tasks WHERE id = $1
RETURNING *;
