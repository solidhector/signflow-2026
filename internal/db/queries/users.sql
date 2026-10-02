-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: GetUserName :one
SELECT full_name FROM users WHERE id = $1;