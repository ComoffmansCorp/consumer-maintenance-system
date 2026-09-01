-- name: UpsertMasterProfile :one
INSERT INTO master_profiles (user_id, city, bio, updated_at)
VALUES (@user_id, @city, @bio, NOW())
ON CONFLICT (user_id) DO UPDATE
SET city = EXCLUDED.city, bio = EXCLUDED.bio, updated_at = NOW()
RETURNING user_id, city, bio, created_at, updated_at;

-- name: GetMasterProfile :one
SELECT user_id, city, bio, created_at, updated_at
FROM master_profiles
WHERE user_id = @user_id
LIMIT 1;

-- name: ReplaceMasterSpecializations :exec
DELETE FROM master_specializations WHERE master_user_id = @master_user_id;

-- name: AddMasterSpecialization :exec
INSERT INTO master_specializations (master_user_id, service_id)
VALUES (@master_user_id, @service_id)
ON CONFLICT DO NOTHING;

-- name: ListMasterSpecializationIDs :many
SELECT service_id
FROM master_specializations
WHERE master_user_id = @master_user_id;

-- name: MasterHasSpecialization :one
SELECT EXISTS (
    SELECT 1 FROM master_specializations
    WHERE master_user_id = @master_user_id AND service_id = @service_id
) AS has_specialization;
