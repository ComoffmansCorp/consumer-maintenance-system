-- name: GetUserByPlatformUsername :one
SELECT id, username, password, full_name, role, tenant_id, created_at
FROM users
WHERE tenant_id IS NULL AND LOWER(username) = LOWER(@username)
LIMIT 1;

-- name: GetUserByTenantCodeAndUsername :one
SELECT u.id, u.username, u.password, u.full_name, u.role, u.tenant_id, u.created_at
FROM users u
JOIN tenants t ON t.id = u.tenant_id
WHERE LOWER(t.code) = LOWER(@tenant_code)
  AND LOWER(u.username) = LOWER(@username)
LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, password, full_name, role, tenant_id, created_at
FROM users
WHERE id = @id
LIMIT 1;

-- name: CountUsersByRole :one
SELECT COUNT(*)::bigint AS count
FROM users
WHERE role = @role;

-- name: CountUsersByTenantID :one
SELECT COUNT(*)::bigint AS count
FROM users
WHERE tenant_id = @tenant_id::bigint;

-- name: ExistsUserInTenant :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE tenant_id = @tenant_id
      AND LOWER(username) = LOWER(@username)
) AS exists;

-- name: ExistsPlatformUser :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE tenant_id IS NULL
      AND LOWER(username) = LOWER(@username)
) AS exists;

-- name: CreateUser :one
INSERT INTO users (username, password, full_name, role, tenant_id)
VALUES (@username, @password, @full_name, @role, @tenant_id)
RETURNING id, username, password, full_name, role, tenant_id, created_at;

-- name: ListUsersByTenant :many
SELECT id, username, password, full_name, role, tenant_id, created_at
FROM users
WHERE tenant_id = @tenant_id::bigint
  AND (@role::text = '' OR role = @role::text)
ORDER BY full_name
LIMIT @page_limit OFFSET @page_offset;

-- name: CountUsersByTenantAndRole :one
SELECT COUNT(*)::bigint AS count
FROM users
WHERE tenant_id = @tenant_id::bigint
  AND (@role::text = '' OR role = @role::text);

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES (@user_id, @token_hash, @expires_at)
RETURNING id, user_id, token_hash, expires_at, created_at, revoked_at;

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
FROM refresh_tokens
WHERE token_hash = @token_hash
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token_hash = @token_hash;

-- name: RevokeUserRefreshTokens :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE user_id = @user_id
  AND revoked_at IS NULL;
