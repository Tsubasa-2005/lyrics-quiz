package infra

import (
	"lyrics-quiz/pkg/infra/rdb"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/taxio/errors"
)

func Initialize(c *gin.Context, userID string, bot *linebot.Client, event *linebot.Event) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	reply := "エラーが発生しました。もう一度やり直してください。"

	err := repo.UpdateQuizManager(c, rdb.UpdateQuizManagerParams{
		TheNumberOfQuestions: 0,
		QuizCount:            1,
		LyricsCount:          1,
		Status:               "not_started",
		Type:                 "",
		UserID:               userID,
	})
	if err != nil {
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
			c.Error(errors.Wrap(err))
		}
	}
	err = repo.DeleteLyrics(c, userID)
	if err != nil {
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
			c.Error(errors.Wrap(err))
		}
	}
	err = repo.DeleteAnswer(c, userID)
	if err != nil {
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
			c.Error(errors.Wrap(err))
		}
	}
	err = repo.DeleteArtist(c, userID)
	if err != nil {
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
			c.Error(errors.Wrap(err))
		}
	}
	err = repo.DeleteChoices(c, userID)
	if err != nil {
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
			c.Error(errors.Wrap(err))
		}
	}
}
