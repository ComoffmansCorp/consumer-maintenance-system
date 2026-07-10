-- name: CreateTask :one
INSERT INTO tasks (type, tenant_id, address_id, status, due_date, assignee_id)
VALUES (@type, @tenant_id, @address_id, 'PENDING', @due_date, @assignee_id)
RETURNING id, type, tenant_id, address_id, status, due_date, assignee_id,
    created_at, updated_at, completed_at, canceled_at, cancel_reason;

-- name: GetTaskByID :one
SELECT id, type, tenant_id, address_id, status, due_date, assignee_id,
    created_at, updated_at, completed_at, canceled_at, cancel_reason
FROM tasks
WHERE id = @id AND tenant_id = @tenant_id
LIMIT 1;

-- name: ListTasks :many
SELECT id, type, tenant_id, address_id, status, due_date, assignee_id,
    created_at, updated_at, completed_at, canceled_at, cancel_reason
FROM tasks
WHERE tenant_id = @tenant_id
  AND (@status::text = '' OR status = @status::text)
  AND (@task_type::text = '' OR type = @task_type::text)
  AND (@assignee_id::bigint = 0 OR assignee_id = @assignee_id::bigint)
ORDER BY
    CASE status WHEN 'PENDING' THEN 0 WHEN 'IN_PROGRESS' THEN 1 ELSE 2 END,
    due_date NULLS LAST,
    created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountTasks :one
SELECT COUNT(*)::bigint AS count
FROM tasks
WHERE tenant_id = @tenant_id
  AND (@status::text = '' OR status = @status::text)
  AND (@task_type::text = '' OR type = @task_type::text)
  AND (@assignee_id::bigint = 0 OR assignee_id = @assignee_id::bigint);

-- name: AssignTask :one
UPDATE tasks
SET assignee_id = @assignee_id, updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id AND status = 'PENDING'
RETURNING id, type, tenant_id, address_id, status, due_date, assignee_id,
    created_at, updated_at, completed_at, canceled_at, cancel_reason;

-- name: StartTask :one
UPDATE tasks
SET status = 'IN_PROGRESS', updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id AND status = 'PENDING'
RETURNING id, type, tenant_id, address_id, status, due_date, assignee_id,
    created_at, updated_at, completed_at, canceled_at, cancel_reason;

-- name: CompleteTask :one
UPDATE tasks
SET status = 'COMPLETED', completed_at = NOW(), updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id AND status = 'IN_PROGRESS'
RETURNING id, type, tenant_id, address_id, status, due_date, assignee_id,
    created_at, updated_at, completed_at, canceled_at, cancel_reason;

-- name: CancelTask :one
UPDATE tasks
SET status = 'CANCELED', canceled_at = NOW(), cancel_reason = @cancel_reason, updated_at = NOW()
WHERE id = @id AND tenant_id = @tenant_id AND status IN ('PENDING', 'IN_PROGRESS')
RETURNING id, type, tenant_id, address_id, status, due_date, assignee_id,
    created_at, updated_at, completed_at, canceled_at, cancel_reason;
