-- name: CreateQuizManager :exec
INSERT INTO quiz_manager (user_id, the_number_of_questions, quiz_count, lyrics_count, status, type) VALUES (?1, ?2, ?3, ?4, ?5, ?6);

-- name: GetQuizManager :one
SELECT * FROM quiz_manager WHERE user_id = ?1;

-- name: UpdateQuizManager :exec
UPDATE quiz_manager SET the_number_of_questions = ?1, quiz_count = ?2, lyrics_count = ?3, status = ?4, type = ?5 WHERE user_id = ?6;