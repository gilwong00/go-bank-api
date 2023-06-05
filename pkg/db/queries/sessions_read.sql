-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1
LIMIT 1;