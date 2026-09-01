-- name: DeleteMasterSpecializations :exec
DELETE FROM master_specializations WHERE master_user_id = @master_user_id;

-- name: AddMasterSpecialization :exec
INSERT INTO master_specializations (master_user_id, service_id)
VALUES (@master_user_id, @service_id)
ON CONFLICT DO NOTHING;

-- name: ListMasterSpecializationIDs :many
SELECT service_id
FROM master_specializations
WHERE master_user_id = @master_user_id
ORDER BY service_id;

-- name: MasterHasSpecialization :one
SELECT EXISTS(
    SELECT 1
    FROM master_specializations
    WHERE master_user_id = @master_user_id AND service_id = @service_id
) AS exists;
