-- name: ListActiveServices :many
SELECT id, category_id, name, description, active, created_at
FROM services
WHERE active = TRUE
  AND (@category_id::bigint = 0 OR category_id = @category_id::bigint)
ORDER BY name;

-- name: GetServiceByID :one
SELECT id, category_id, name, description, active, created_at
FROM services
WHERE id = @id
LIMIT 1;

-- name: ListServicesByIDs :many
SELECT id, category_id, name, description, active, created_at
FROM services
WHERE id = ANY(@ids::bigint[]);

-- name: CreateService :one
INSERT INTO services (category_id, name, description)
VALUES (@category_id, @name, @description)
RETURNING id, category_id, name, description, active, created_at;
