-- name: AddFavorite :exec
INSERT INTO favorites (client_id, master_id)
VALUES (@client_id, @master_id)
ON CONFLICT DO NOTHING;

-- name: RemoveFavorite :exec
DELETE FROM favorites WHERE client_id = @client_id AND master_id = @master_id;

-- name: ListFavoritesByClient :many
SELECT client_id, master_id, created_at
FROM favorites
WHERE client_id = @client_id
ORDER BY created_at DESC;
