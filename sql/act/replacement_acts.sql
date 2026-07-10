-- name: CreateReplacementAct :one
INSERT INTO replacement_acts (
    task_id, tenant_id, address_id, account_number, installation_date,
    old_brand, old_serial_number, old_readings, new_brand, new_serial_number, new_readings
)
VALUES (
    @task_id, @tenant_id, @address_id, @account_number, @installation_date,
    @old_brand, @old_serial_number, @old_readings, @new_brand, @new_serial_number, @new_readings
)
RETURNING id, task_id, tenant_id, address_id, account_number, installation_date,
    old_brand, old_serial_number, old_readings, new_brand, new_serial_number, new_readings,
    created_at, updated_at;

-- name: GetReplacementActByID :one
SELECT id, task_id, tenant_id, address_id, account_number, installation_date,
    old_brand, old_serial_number, old_readings, new_brand, new_serial_number, new_readings,
    created_at, updated_at
FROM replacement_acts
WHERE id = @id AND tenant_id = @tenant_id
LIMIT 1;

-- name: GetReplacementActByTaskID :one
SELECT id, task_id, tenant_id, address_id, account_number, installation_date,
    old_brand, old_serial_number, old_readings, new_brand, new_serial_number, new_readings,
    created_at, updated_at
FROM replacement_acts
WHERE task_id = @task_id AND tenant_id = @tenant_id
LIMIT 1;

-- name: ExistsReplacementActByTaskID :one
SELECT EXISTS(SELECT 1 FROM replacement_acts WHERE task_id = @task_id) AS exists;

-- name: UpdateReplacementAct :one
UPDATE replacement_acts
SET account_number = @account_number,
    installation_date = @installation_date,
    old_brand = @old_brand,
    old_serial_number = @old_serial_number,
    old_readings = @old_readings,
    new_brand = @new_brand,
    new_serial_number = @new_serial_number,
    new_readings = @new_readings,
    updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, task_id, tenant_id, address_id, account_number, installation_date,
    old_brand, old_serial_number, old_readings, new_brand, new_serial_number, new_readings,
    created_at, updated_at;
