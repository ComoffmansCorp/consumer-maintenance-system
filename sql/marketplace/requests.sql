-- name: CreateServiceRequest :one
INSERT INTO service_requests (client_id, service_id, description, address_text, latitude, longitude, status)
VALUES (@client_id, @service_id, @description, @address_text, @latitude, @longitude, 'OPEN')
RETURNING id, client_id, service_id, description, address_text, latitude, longitude, status, master_id,
    created_at, updated_at, claimed_at, completed_at, canceled_at, cancel_reason;

-- name: GetServiceRequestByID :one
SELECT id, client_id, service_id, description, address_text, latitude, longitude, status, master_id,
    created_at, updated_at, claimed_at, completed_at, canceled_at, cancel_reason
FROM service_requests
WHERE id = @id
LIMIT 1;

-- name: ListRequestsByClient :many
SELECT id, client_id, service_id, description, address_text, latitude, longitude, status, master_id,
    created_at, updated_at, claimed_at, completed_at, canceled_at, cancel_reason
FROM service_requests
WHERE client_id = @client_id
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountRequestsByClient :one
SELECT COUNT(*)::bigint AS count
FROM service_requests
WHERE client_id = @client_id;

-- name: ListRequestsByMaster :many
SELECT id, client_id, service_id, description, address_text, latitude, longitude, status, master_id,
    created_at, updated_at, claimed_at, completed_at, canceled_at, cancel_reason
FROM service_requests
WHERE master_id = @master_id
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountRequestsByMaster :one
SELECT COUNT(*)::bigint AS count
FROM service_requests
WHERE master_id = @master_id;

-- name: ListOpenRequestsForMaster :many
-- Pool of unclaimed requests, restricted to services the master is
-- specialized in -- the hard server-side filter, not just a UI convenience.
SELECT sr.id, sr.client_id, sr.service_id, sr.description, sr.address_text, sr.latitude, sr.longitude,
    sr.status, sr.master_id, sr.created_at, sr.updated_at, sr.claimed_at, sr.completed_at,
    sr.canceled_at, sr.cancel_reason
FROM service_requests sr
JOIN master_specializations ms
    ON ms.service_id = sr.service_id AND ms.master_user_id = @master_user_id
WHERE sr.status = 'OPEN'
ORDER BY sr.created_at
LIMIT @page_limit OFFSET @page_offset;

-- name: CountOpenRequestsForMaster :one
SELECT COUNT(*)::bigint AS count
FROM service_requests sr
JOIN master_specializations ms
    ON ms.service_id = sr.service_id AND ms.master_user_id = @master_user_id
WHERE sr.status = 'OPEN';

-- name: ClaimServiceRequest :one
UPDATE service_requests
SET status = 'IN_PROGRESS', master_id = @master_id, claimed_at = NOW(), updated_at = NOW()
WHERE id = @id AND status = 'OPEN'
RETURNING id, client_id, service_id, description, address_text, latitude, longitude, status, master_id,
    created_at, updated_at, claimed_at, completed_at, canceled_at, cancel_reason;

-- name: CompleteServiceRequest :one
UPDATE service_requests
SET status = 'COMPLETED', completed_at = NOW(), updated_at = NOW()
WHERE id = @id AND status = 'IN_PROGRESS'
RETURNING id, client_id, service_id, description, address_text, latitude, longitude, status, master_id,
    created_at, updated_at, claimed_at, completed_at, canceled_at, cancel_reason;

-- name: CancelServiceRequest :one
UPDATE service_requests
SET status = 'CANCELED', canceled_at = NOW(), cancel_reason = @cancel_reason, updated_at = NOW()
WHERE id = @id AND status IN ('OPEN', 'IN_PROGRESS')
RETURNING id, client_id, service_id, description, address_text, latitude, longitude, status, master_id,
    created_at, updated_at, claimed_at, completed_at, canceled_at, cancel_reason;
