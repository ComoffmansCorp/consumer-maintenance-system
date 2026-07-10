-- name: CreateMeter :one
INSERT INTO meters (type, serial_number, manufacture_year, verification_date, seal_state, transformation_ratio, inspection_act_id)
VALUES (@type, @serial_number, @manufacture_year, @verification_date, @seal_state, @transformation_ratio, @inspection_act_id)
RETURNING id, type, serial_number, manufacture_year, verification_date, seal_state, transformation_ratio, inspection_act_id, created_at;

-- name: GetMeterByID :one
SELECT id, type, serial_number, manufacture_year, verification_date, seal_state, transformation_ratio, inspection_act_id, created_at
FROM meters
WHERE id = @id AND inspection_act_id = @inspection_act_id
LIMIT 1;

-- name: ListMetersByAct :many
SELECT id, type, serial_number, manufacture_year, verification_date, seal_state, transformation_ratio, inspection_act_id, created_at
FROM meters
WHERE inspection_act_id = @inspection_act_id
ORDER BY created_at;

-- name: UpdateMeter :one
UPDATE meters
SET type = @type,
    serial_number = @serial_number,
    manufacture_year = @manufacture_year,
    verification_date = @verification_date,
    seal_state = @seal_state,
    transformation_ratio = @transformation_ratio
WHERE id = @id AND inspection_act_id = @inspection_act_id
RETURNING id, type, serial_number, manufacture_year, verification_date, seal_state, transformation_ratio, inspection_act_id, created_at;

-- name: DeleteMeter :exec
DELETE FROM meters
WHERE id = @id AND inspection_act_id = @inspection_act_id;
