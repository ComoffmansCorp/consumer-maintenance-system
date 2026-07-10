-- name: CreateAddress :one
INSERT INTO addresses (street, house, building, apartment, tenant_id, consumer_id)
VALUES (@street, @house, @building, @apartment, @tenant_id, @consumer_id)
RETURNING id, street, house, building, apartment, tenant_id, consumer_id, created_at, updated_at;

-- name: GetAddressByID :one
SELECT id, street, house, building, apartment, tenant_id, consumer_id, created_at, updated_at
FROM addresses
WHERE id = @id AND tenant_id = @tenant_id
LIMIT 1;

-- name: ExistsAddressByID :one
SELECT EXISTS(
    SELECT 1 FROM addresses WHERE id = @id AND tenant_id = @tenant_id
) AS exists;

-- name: ListAddresses :many
SELECT id, street, house, building, apartment, tenant_id, consumer_id, created_at, updated_at
FROM addresses
WHERE tenant_id = @tenant_id
  AND (@search::text = '' OR street ILIKE '%' || @search::text || '%' OR house ILIKE '%' || @search::text || '%')
  AND (@consumer_id::bigint = 0 OR consumer_id = @consumer_id::bigint)
ORDER BY street, house
LIMIT @page_limit OFFSET @page_offset;

-- name: CountAddresses :one
SELECT COUNT(*)::bigint AS count
FROM addresses
WHERE tenant_id = @tenant_id
  AND (@search::text = '' OR street ILIKE '%' || @search::text || '%' OR house ILIKE '%' || @search::text || '%')
  AND (@consumer_id::bigint = 0 OR consumer_id = @consumer_id::bigint);

-- name: UpdateAddress :one
UPDATE addresses
SET street = @street,
    house = @house,
    building = @building,
    apartment = @apartment,
    consumer_id = @consumer_id,
    updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, street, house, building, apartment, tenant_id, consumer_id, created_at, updated_at;
