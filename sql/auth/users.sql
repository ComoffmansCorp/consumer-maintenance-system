-- name: CreateUser :one
INSERT INTO users (username, password_hash, full_name, role)
VALUES (@username, @password_hash, @full_name, @role)
RETURNING id, username, password_hash, full_name, role, created_at, updated_at;

-- name: GetUserByUsername :one
SELECT id, username, password_hash, full_name, role, created_at, updated_at
FROM users
WHERE LOWER(username) = LOWER(@username)
LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, password_hash, full_name, role, created_at, updated_at
FROM users
WHERE id = @id
LIMIT 1;

-- name: CountUsersByRole :one
SELECT COUNT(*)::bigint AS count
FROM users
WHERE role = @role;

-- name: ExistsUserByUsername :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE LOWER(username) = LOWER(@username)
) AS exists;
