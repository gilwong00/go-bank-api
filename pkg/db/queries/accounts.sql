-- name: CreateAccount :one
INSERT INTO accounts (
	owner,
	balance,
	currency
) VALUES (
	$1, $2, $3
) RETURNING *;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE from accounts
WHERE id = $1;

-- name: GetAccountById :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccounts :many
SELECT * from accounts
ORDER BY name;