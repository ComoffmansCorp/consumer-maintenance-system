-- name: CreateStatusHistory :exec
INSERT INTO request_status_history (request_id, from_status, to_status, changed_by, comment)
VALUES (@request_id, @from_status, @to_status, @changed_by, @comment);

-- name: ListStatusHistoryByRequest :many
SELECT id, request_id, from_status, to_status, changed_by, comment, created_at
FROM request_status_history
WHERE request_id = @request_id
ORDER BY created_at;
