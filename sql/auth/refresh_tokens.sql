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
