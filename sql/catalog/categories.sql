-- name: CreateCategory :one
INSERT INTO service_categories (parent_category_id, name)
VALUES (@parent_category_id, @name)
RETURNING id, parent_category_id, name, active, created_at, updated_at;

-- name: GetCategoryByID :one
SELECT id, parent_category_id, name, active, created_at, updated_at
FROM service_categories
WHERE id = @id
LIMIT 1;

-- name: UpdateCategory :one
UPDATE service_categories
SET name = @name, active = @active, updated_at = NOW()
WHERE id = @id
RETURNING id, parent_category_id, name, active, created_at, updated_at;

-- name: ListActiveCategories :many
-- Every active category, top-level and subcategories alike -- the service
-- layer groups children under their parent to build the nested tree.
SELECT id, parent_category_id, name, active, created_at, updated_at
FROM service_categories
WHERE active = true
ORDER BY parent_category_id NULLS FIRST, name;
