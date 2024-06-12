-- name: CreateLyrics :exec
INSERT INTO lyrics (quiz_manager_id, question_number, count, lyrics) VALUES (?1, ?2, ?3, ?4);

-- name: GetLyrics :one
SELECT * FROM lyrics WHERE quiz_manager_id = ?1 AND question_number = ?2 AND count = ?3;

-- name: DeleteLyrics :exec
DELETE FROM lyrics WHERE quiz_manager_id = ?1;