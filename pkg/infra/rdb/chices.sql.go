// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: chices.sql

package rdb

import (
	"context"
)

const createChoices = `-- name: CreateChoices :exec
INSERT INTO choices (quiz_manager_id, question_number, choice1, choice2, choice3, choice4) VALUES (?1, ?2, ?3, ?4, ?5, ?6)
`

type CreateChoicesParams struct {
	QuizManagerID  string
	QuestionNumber int64
	Choice1        string
	Choice2        string
	Choice3        string
	Choice4        string
}

func (q *Queries) CreateChoices(ctx context.Context, arg CreateChoicesParams) error {
	_, err := q.db.ExecContext(ctx, createChoices,
		arg.QuizManagerID,
		arg.QuestionNumber,
		arg.Choice1,
		arg.Choice2,
		arg.Choice3,
		arg.Choice4,
	)
	return err
}

const deleteChoices = `-- name: DeleteChoices :exec
DELETE FROM choices WHERE quiz_manager_id = ?1
`

func (q *Queries) DeleteChoices(ctx context.Context, quizManagerID string) error {
	_, err := q.db.ExecContext(ctx, deleteChoices, quizManagerID)
	return err
}

const getChoices = `-- name: GetChoices :one
SELECT id, quiz_manager_id, question_number, choice1, choice2, choice3, choice4 FROM choices WHERE quiz_manager_id = ?1 AND question_number = ?2
`

type GetChoicesParams struct {
	QuizManagerID  string
	QuestionNumber int64
}

func (q *Queries) GetChoices(ctx context.Context, arg GetChoicesParams) (Choice, error) {
	row := q.db.QueryRowContext(ctx, getChoices, arg.QuizManagerID, arg.QuestionNumber)
	var i Choice
	err := row.Scan(
		&i.ID,
		&i.QuizManagerID,
		&i.QuestionNumber,
		&i.Choice1,
		&i.Choice2,
		&i.Choice3,
		&i.Choice4,
	)
	return i, err
}