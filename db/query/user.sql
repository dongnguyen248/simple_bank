-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserForUdpate :one
SELECT * FROM users
WHERE username = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUser :many
SELECT * FROM users
ORDER BY username 
LIMIT $1
OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
    )
VALUES
  (
    $1,$2,$3,$4
  ) 
  RETURNING *;
-- name: ChangePassword :one
UPDATE users
SET hashed_password = $1 ,
    password_changed_at = now()
WHERE username = sqlc.arg(username)
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;