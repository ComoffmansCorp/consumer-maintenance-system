-- name: CreateInspectionAct :one
INSERT INTO inspection_acts (task_id, tenant_id, address_id, inspection_date, consumer_id, inspection_type, notes)
VALUES (@task_id, @tenant_id, @address_id, @inspection_date, @consumer_id, @inspection_type, @notes)
RETURNING id, task_id, tenant_id, address_id, inspection_date, consumer_id, inspection_type, notes, created_at, updated_at;

-- name: GetInspectionActByID :one
SELECT id, task_id, tenant_id, address_id, inspection_date, consumer_id, inspection_type, notes, created_at, updated_at
FROM inspection_acts
WHERE id = @id AND tenant_id = @tenant_id
LIMIT 1;

-- name: GetInspectionActByTaskID :one
SELECT id, task_id, tenant_id, address_id, inspection_date, consumer_id, inspection_type, notes, created_at, updated_at
FROM inspection_acts
WHERE task_id = @task_id AND tenant_id = @tenant_id
LIMIT 1;

-- name: ExistsInspectionActByTaskID :one
SELECT EXISTS(SELECT 1 FROM inspection_acts WHERE task_id = @task_id) AS exists;

-- name: UpdateInspectionAct :one
UPDATE inspection_acts
SET inspection_date = @inspection_date,
    consumer_id = @consumer_id,
    inspection_type = @inspection_type,
    notes = @notes,
    updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, task_id, tenant_id, address_id, inspection_date, consumer_id, inspection_type, notes, created_at, updated_at;
