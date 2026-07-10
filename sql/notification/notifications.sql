-- name: CreateNotification :one
INSERT INTO notifications (tenant_id, user_id, type, title, message, payload)
VALUES (@tenant_id, @user_id, @type, @title, @message, @payload)
RETURNING id, tenant_id, user_id, type, title, message, payload, read_at, created_at;

-- name: ListNotificationsByUser :many
SELECT id, tenant_id, user_id, type, title, message, payload, read_at, created_at
FROM notifications
WHERE tenant_id = @tenant_id AND user_id = @user_id
  AND (@unread_only::bool = FALSE OR read_at IS NULL)
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountNotificationsByUser :one
SELECT COUNT(*)::bigint AS count
FROM notifications
WHERE tenant_id = @tenant_id AND user_id = @user_id
  AND (@unread_only::bool = FALSE OR read_at IS NULL);

-- name: CountUnreadByUser :one
SELECT COUNT(*)::bigint AS count
FROM notifications
WHERE tenant_id = @tenant_id AND user_id = @user_id AND read_at IS NULL;

-- name: MarkNotificationRead :one
UPDATE notifications
SET read_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id AND user_id = @user_id AND read_at IS NULL
RETURNING id, tenant_id, user_id, type, title, message, payload, read_at, created_at;

-- name: MarkAllNotificationsRead :exec
UPDATE notifications
SET read_at = NOW()
WHERE tenant_id = @tenant_id AND user_id = @user_id AND read_at IS NULL;
