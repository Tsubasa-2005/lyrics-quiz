-- name: CreateArtist :exec
INSERT INTO artist (quiz_manager_id, artist) VALUES (?1, ?2);

-- name: GetArtist :one
SELECT * FROM artist WHERE quiz_manager_id = ?1;

-- name: DeleteArtist :exec
DELETE FROM artist WHERE quiz_manager_id = ?1;

-- name: UpdateArtist :exec
UPDATE artist SET artist = ?2 WHERE quiz_manager_id = ?1 returning *;