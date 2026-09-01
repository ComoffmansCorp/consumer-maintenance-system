-- name: CreateReview :one
INSERT INTO reviews (request_id, client_id, master_id, rating, comment, hidden)
VALUES (@request_id, @client_id, @master_id, @rating, @comment, false)
RETURNING id, request_id, client_id, master_id, rating, comment, hidden, created_at;

-- name: GetReviewByID :one
SELECT id, request_id, client_id, master_id, rating, comment, hidden, created_at
FROM reviews
WHERE id = @id
LIMIT 1;

-- name: GetReviewByRequestID :one
SELECT id, request_id, client_id, master_id, rating, comment, hidden, created_at
FROM reviews
WHERE request_id = @request_id
LIMIT 1;

-- name: ListVisibleReviewsByMaster :many
SELECT id, request_id, client_id, master_id, rating, comment, hidden, created_at
FROM reviews
WHERE master_id = @master_id AND hidden = false
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountVisibleReviewsByMaster :one
SELECT COUNT(*)::bigint AS count
FROM reviews
WHERE master_id = @master_id AND hidden = false;

-- name: HideReview :one
UPDATE reviews
SET hidden = true
WHERE id = @id
RETURNING id, request_id, client_id, master_id, rating, comment, hidden, created_at;
