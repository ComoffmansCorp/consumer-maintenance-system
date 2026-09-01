-- name: CreateOffer :one
INSERT INTO request_offers (request_id, master_id, price, comment, status)
VALUES (@request_id, @master_id, @price, @comment, 'PENDING')
RETURNING id, request_id, master_id, price, comment, status, created_at, updated_at;

-- name: GetOfferByID :one
SELECT id, request_id, master_id, price, comment, status, created_at, updated_at
FROM request_offers
WHERE id = @id
LIMIT 1;

-- name: ListOffersByRequest :many
SELECT id, request_id, master_id, price, comment, status, created_at, updated_at
FROM request_offers
WHERE request_id = @request_id
ORDER BY created_at;

-- name: AcceptOffer :one
UPDATE request_offers
SET status = 'ACCEPTED', updated_at = NOW()
WHERE id = @id AND status = 'PENDING'
RETURNING id, request_id, master_id, price, comment, status, created_at, updated_at;

-- name: RejectOtherOffers :exec
UPDATE request_offers
SET status = 'REJECTED', updated_at = NOW()
WHERE request_id = @request_id AND id != @accepted_offer_id AND status = 'PENDING';
