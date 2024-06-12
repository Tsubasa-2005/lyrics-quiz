package handler

import (
	"lyrics-quiz/pkg/infra/rdb"
	"lyrics-quiz/pkg/message"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func UnknownQuestion(c *gin.Context, quizManager rdb.QuizManager, bot *linebot.Client, event *linebot.Event) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	answer, err := repo.GetAnswer(c, rdb.GetAnswerParams{
		QuizManagerID:  quizManager.UserID,
		QuestionNumber: quizManager.QuizCount,
	})
	if err != nil {
		message.Error(c, bot, event)
		return
	}
	quizManager.LyricsCount = 1
	quizManager.QuizCount++
	err = repo.UpdateQuizManager(c, rdb.UpdateQuizManagerParams{
		TheNumberOfQuestions: quizManager.TheNumberOfQuestions,
		QuizCount:            quizManager.QuizCount,
		LyricsCount:          quizManager.LyricsCount,
		Status:               quizManager.Status,
		Type:                 quizManager.Type,
		UserID:               quizManager.UserID,
	})
	if err != nil {
		message.Error(c, bot, event)
		return
	}
	message.Unknown(c, bot, event, answer.MusicName)
}
