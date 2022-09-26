-- name: DeleteAccount :exec
DELETE from accounts
WHERE id = $1;