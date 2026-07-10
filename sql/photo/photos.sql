-- name: CreatePhoto :one
INSERT INTO photos (
    filename, note, tenant_id, inspection_act_id, replacement_act_id,
    original_filename, content_type, size_bytes, uploaded_by
)
VALUES (
    @filename, @note, @tenant_id, @inspection_act_id, @replacement_act_id,
    @original_filename, @content_type, @size_bytes, @uploaded_by
)
RETURNING id, filename, note, tenant_id, inspection_act_id, replacement_act_id,
    original_filename, content_type, size_bytes, uploaded_by, created_at;

-- name: GetPhotoByID :one
SELECT id, filename, note, tenant_id, inspection_act_id, replacement_act_id,
    original_filename, content_type, size_bytes, uploaded_by, created_at
FROM photos
WHERE id = @id AND tenant_id = @tenant_id
LIMIT 1;

-- name: ListPhotosByInspectionAct :many
SELECT id, filename, note, tenant_id, inspection_act_id, replacement_act_id,
    original_filename, content_type, size_bytes, uploaded_by, created_at
FROM photos
WHERE inspection_act_id = @inspection_act_id
ORDER BY created_at;

-- name: ListPhotosByReplacementAct :many
SELECT id, filename, note, tenant_id, inspection_act_id, replacement_act_id,
    original_filename, content_type, size_bytes, uploaded_by, created_at
FROM photos
WHERE replacement_act_id = @replacement_act_id
ORDER BY created_at;

-- name: DeletePhoto :one
DELETE FROM photos
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, filename, note, tenant_id, inspection_act_id, replacement_act_id,
    original_filename, content_type, size_bytes, uploaded_by, created_at;
