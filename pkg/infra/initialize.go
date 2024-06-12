package infra

import (
	"lyrics-quiz/pkg/infra/rdb"
	"lyrics-quiz/pkg/message"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func Initialize(c *gin.Context, userID string, bot *linebot.Client, event *linebot.Event) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	err := repo.UpdateQuizManager(c, rdb.UpdateQuizManagerParams{
		TheNumberOfQuestions: 0,
		QuizCount:            1,
		LyricsCount:          1,
		Status:               "not_started",
		Type:                 "",
		UserID:               userID,
	})
	if err != nil {
		message.Error(c, bot, event)
	}
	err = repo.DeleteLyrics(c, userID)
	if err != nil {
		message.Error(c, bot, event)
	}
	err = repo.DeleteAnswer(c, userID)
	if err != nil {
		message.Error(c, bot, event)
	}
	err = repo.DeleteArtist(c, userID)
	if err != nil {
		message.Error(c, bot, event)
	}
	err = repo.DeleteChoices(c, userID)
	if err != nil {
		message.Error(c, bot, event)
	}
}
