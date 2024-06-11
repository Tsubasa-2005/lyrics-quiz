-- name: CreateAnswer :exec
INSERT INTO answer (quiz_manager_id, question_number, music_name) VALUES (?1, ?2, ?3);

-- name: GetAnswer :one
SELECT * FROM answer WHERE quiz_manager_id = ?1 AND question_number = ?2;

-- name: DeleteAnswer :exec
DELETE FROM answer WHERE quiz_manager_id = ?1;