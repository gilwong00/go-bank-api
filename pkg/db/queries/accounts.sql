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

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
-- tells postgres that we dont update the key of id column of the account table to avoid deadlock
FOR NO KEY UPDATE;

-- name: GetAccounts :many
SELECT * from accounts
ORDER BY name;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

