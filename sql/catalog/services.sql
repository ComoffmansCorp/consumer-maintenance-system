-- name: CreateService :one
INSERT INTO services (category_id, name, description, price_from, price_to, unit)
VALUES (@category_id, @name, @description, @price_from, @price_to, @unit)
RETURNING id, category_id, name, description, price_from, price_to, unit, active, created_at, updated_at;

-- name: GetServiceByID :one
SELECT id, category_id, name, description, price_from, price_to, unit, active, created_at, updated_at
FROM services
WHERE id = @id
LIMIT 1;

-- name: UpdateService :one
UPDATE services
SET name = @name, description = @description, price_from = @price_from,
    price_to = @price_to, unit = @unit, active = @active, updated_at = NOW()
WHERE id = @id
RETURNING id, category_id, name, description, price_from, price_to, unit, active, created_at, updated_at;

-- name: ListActiveServices :many
-- category_id = 0 means "no filter": every active service across the catalog.
SELECT id, category_id, name, description, price_from, price_to, unit, active, created_at, updated_at
FROM services
WHERE active = true
  AND (@category_id::bigint = 0 OR category_id = @category_id::bigint)
ORDER BY name;
