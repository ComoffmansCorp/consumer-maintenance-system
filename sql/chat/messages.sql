-- name: CreateMessage :one
INSERT INTO messages (request_id, sender_id, text)
VALUES (@request_id, @sender_id, @text)
RETURNING id, request_id, sender_id, text, created_at, read_at;

-- name: ListMessagesByRequestSince :many
-- Cursor pagination: since_id = 0 returns from the beginning of the thread.
SELECT id, request_id, sender_id, text, created_at, read_at
FROM messages
WHERE request_id = @request_id AND id > @since_id::bigint
ORDER BY id
LIMIT @page_limit;
