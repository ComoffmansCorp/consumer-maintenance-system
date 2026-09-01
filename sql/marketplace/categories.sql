-- name: ListActiveCategories :many
SELECT id, name, active, created_at
FROM service_categories
WHERE active = TRUE
ORDER BY name;

-- name: GetCategoryByID :one
SELECT id, name, active, created_at
FROM service_categories
WHERE id = @id
LIMIT 1;

-- name: CreateCategory :one
INSERT INTO service_categories (name)
VALUES (@name)
RETURNING id, name, active, created_at;
