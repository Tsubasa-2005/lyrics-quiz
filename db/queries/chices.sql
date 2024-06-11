-- name: CreateChoices :exec
INSERT INTO choices (quiz_manager_id, question_number, choice1, choice2, choice3, choice4) VALUES (?1, ?2, ?3, ?4, ?5, ?6);

-- name: GetChoices :one
SELECT * FROM choices WHERE quiz_manager_id = ?1 AND question_number = ?2;

-- name: DeleteChoices :exec
DELETE FROM choices WHERE quiz_manager_id = ?1;