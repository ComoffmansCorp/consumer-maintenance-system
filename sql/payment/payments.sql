-- name: CreatePayment :one
INSERT INTO payments (request_id, amount, platform_fee, status)
VALUES (@request_id, @amount, @platform_fee, 'HELD')
RETURNING id, request_id, amount, platform_fee, status, created_at, updated_at;

-- name: GetPaymentByRequestID :one
SELECT id, request_id, amount, platform_fee, status, created_at, updated_at
FROM payments
WHERE request_id = @request_id
LIMIT 1;

-- name: ReleasePayment :one
UPDATE payments
SET status = 'RELEASED', updated_at = NOW()
WHERE request_id = @request_id AND status = 'HELD'
RETURNING id, request_id, amount, platform_fee, status, created_at, updated_at;

-- name: RefundPayment :one
UPDATE payments
SET status = 'REFUNDED', updated_at = NOW()
WHERE request_id = @request_id AND status = 'HELD'
RETURNING id, request_id, amount, platform_fee, status, created_at, updated_at;

-- name: ListPaymentsAdmin :many
SELECT id, request_id, amount, platform_fee, status, created_at, updated_at
FROM payments
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountPaymentsAdmin :one
SELECT COUNT(*)::bigint AS count
FROM payments;
