package handler

import (
	"lyrics-quiz/pkg/infra"
	"lyrics-quiz/pkg/infra/rdb"
	reply "lyrics-quiz/pkg/message"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func CheckAnswer(c *gin.Context, quizManager rdb.QuizManager, message string, bot *linebot.Client, event *linebot.Event) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	answer, err := repo.GetAnswer(c, rdb.GetAnswerParams{
		QuizManagerID:  quizManager.UserID,
		QuestionNumber: quizManager.QuizCount,
	})
	if err != nil {
		reply.Error(c, bot, event)
		return
	}
	if answer.MusicName == message {
		quizManager.LyricsCount = 1
		quizManager.QuizCount++
		if quizManager.QuizCount > quizManager.TheNumberOfQuestions {
			infra.Initialize(c, quizManager.UserID, bot, event)
			reply.EndQuiz(c, bot, event)
			return
		} else {
			err = repo.UpdateQuizManager(c, rdb.UpdateQuizManagerParams{
				TheNumberOfQuestions: quizManager.TheNumberOfQuestions,
				QuizCount:            quizManager.QuizCount,
				LyricsCount:          quizManager.LyricsCount,
				Status:               quizManager.Status,
				Type:                 quizManager.Type,
				UserID:               quizManager.UserID,
			})
			if err != nil {
				reply.Error(c, bot, event)
				return
			}
			reply.CorrectAnswer(c, bot, event)
			return
		}
	} else {
		reply.IncorrectAnswer(c, bot, event)
		return
	}
}
