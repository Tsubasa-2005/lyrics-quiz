package handler

import (
	"lyrics-quiz/pkg/infra/rdb"
	"lyrics-quiz/pkg/message"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func NextQuiz(c *gin.Context, quizManager rdb.QuizManager, bot *linebot.Client, event *linebot.Event) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	quizManager.LyricsCount = 1
	message.ProblemStatement(c, bot, event, quizManager, repo)
}
