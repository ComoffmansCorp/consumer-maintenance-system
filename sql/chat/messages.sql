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

-- name: MarkMessagesRead :exec
-- Marks every message the caller didn't send in this thread as read --
-- only the other participant's messages count as "unread" for the caller,
-- never their own.
UPDATE messages
SET read_at = NOW()
WHERE request_id = @request_id AND sender_id != @reader_id AND read_at IS NULL;

-- name: CountUnreadRequestIDs :one
-- Conversation-level count (how many threads have at least one unread
-- message), not a raw message count -- that's what a header badge should
-- show ("3 заявки с новыми сообщениями"), not "17 unread messages" which
-- reads oddly when most of that is one long back-and-forth.
SELECT COUNT(DISTINCT request_id)::bigint AS count
FROM messages
WHERE request_id = ANY(@request_ids::bigint[]) AND sender_id != @reader_id AND read_at IS NULL;
