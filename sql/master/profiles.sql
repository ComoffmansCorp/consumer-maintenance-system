-- name: UpsertMasterProfile :one
INSERT INTO master_profiles (user_id, city, bio)
VALUES (@user_id, @city, @bio)
ON CONFLICT (user_id) DO UPDATE SET city = excluded.city, bio = excluded.bio, updated_at = NOW()
RETURNING user_id, city, bio, rating_avg, rating_count, created_at, updated_at;

-- name: GetMasterProfile :one
SELECT user_id, city, bio, rating_avg, rating_count, created_at, updated_at
FROM master_profiles
WHERE user_id = @user_id
LIMIT 1;

-- name: ListMasterProfiles :many
SELECT user_id, city, bio, rating_avg, rating_count, created_at, updated_at
FROM master_profiles
ORDER BY rating_avg DESC, created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountMasterProfiles :one
SELECT COUNT(*)::bigint AS count
FROM master_profiles;

-- name: RecordMasterReview :one
-- Atomic running-average update: new_avg = (avg*count + rating) / (count+1).
-- Done in SQL rather than read-modify-write in Go to avoid a lost update
-- race between two reviews landing concurrently for the same master.
UPDATE master_profiles
SET rating_avg = ROUND(((rating_avg * rating_count) + @rating::numeric) / (rating_count + 1), 2),
    rating_count = rating_count + 1,
    updated_at = NOW()
WHERE user_id = @user_id
RETURNING user_id, city, bio, rating_avg, rating_count, created_at, updated_at;
