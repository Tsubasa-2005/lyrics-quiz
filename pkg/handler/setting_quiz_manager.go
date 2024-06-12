package handler

import (
	"lyrics-quiz/pkg/infra/rdb"
	"lyrics-quiz/pkg/message"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func CheckUser(c *gin.Context, userID string, bot *linebot.Client, event *linebot.Event) rdb.QuizManager {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	quizManager, err := repo.GetQuizManager(c, userID)

	if err != nil {
		err = repo.CreateQuizManager(c, rdb.CreateQuizManagerParams{
			UserID:               userID,
			TheNumberOfQuestions: 0,
			QuizCount:            1,
			LyricsCount:          1,
			Status:               "not_started",
			Type:                 "",
		})
		if err != nil {
			message.Error(c, bot, event)
			return rdb.QuizManager{}
		}
		err = repo.CreateArtist(c, rdb.CreateArtistParams{
			QuizManagerID: userID,
			Artist:        "",
		})
		if err != nil {
			message.Error(c, bot, event)
			return rdb.QuizManager{}
		}
		message.FinishedInitialize(c, bot, event)
	}

	return quizManager
}
