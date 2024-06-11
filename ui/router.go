package ui

import (
	"lyrics-quiz/pkg/handler"
	"lyrics-quiz/pkg/infra"
	replyMessage "lyrics-quiz/pkg/message"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/taxio/errors"
)

var bot *linebot.Client

func DBContext() gin.HandlerFunc {
	dbConn := infra.ConnectDB()
	return func(c *gin.Context) {
		c.Set("db", dbConn)
	}
}

func MessageInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if errors.Is(err, linebot.ErrInvalidSignature) {
				c.AbortWithStatus(http.StatusBadRequest)
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					// 特定のメッセージを受け取った場合、処理を中断して終了メッセージを返す
					if message.Text == "e" || message.Text == "E" {
						infra.Initialize(c, event.Source.UserID, bot, event)
						replyMessage.FinishedQuiz(c, bot, event)
						c.Abort()
						return
					}
				}
			}
		}
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	var (
		channelSecret      = os.Getenv("CHANNEL_SECRET")
		channelAccessToken = os.Getenv("CHANNEL_ACCESS_TOKEN")
	)

	var err error
	bot, err = linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(DBContext())

	// ping
	r.GET("/ping/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// db-ping
	r.GET("/db-ping/", func(ctx *gin.Context) {
		infra.CheckConnectDB()
		ctx.JSON(200, gin.H{
			"message": "db-pong",
		})
	})

	r.POST("/webhook", func(ctx *gin.Context) {
		events, err := bot.ParseRequest(ctx.Request)
		if err != nil {
			if errors.Is(err, linebot.ErrInvalidSignature) {
				ctx.AbortWithStatus(http.StatusBadRequest)
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		for _, event := range events {
			UserID := event.Source.UserID
			quizManager := handler.CheckUser(ctx, UserID, bot, event)
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if message.Text == "e" || message.Text == "E" {
						infra.Initialize(ctx, event.Source.UserID, bot, event)
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("クイズを終了しました。")).Do(); err != nil {
							ctx.Error(errors.Wrap(err))
							return
						}
						return
					} else {
						if quizManager.Status == "not_started" {
							if message.Text == "s" || message.Text == "S" {
								handler.AskQuizType(ctx, quizManager, bot, event)
							} else {
								replyMessage.Parroting(ctx, bot, event, message)
							}
						} else if quizManager.Status == "quiz_type" {
							quizType := message.Text
							if quizType == "normal" || quizType == "hard" {
								handler.InputQuizTypeAndAskTheNumberOfQuestions(ctx, quizManager, quizType, bot, event)
							} else {
								title := "選択肢から選んでください。"
								replyMessage.AskQuizTypeMessage(ctx, bot, event, title)
							}
						} else if quizManager.Status == "the_number_of_questions" {
							numberOfQuestions := message.Text
							intNumberOfQuestions, _ := strconv.Atoi(numberOfQuestions)
							if numberOfQuestions == "1" || numberOfQuestions == "5" || numberOfQuestions == "10" {
								handler.InputTheNumberOfQuestionsAndAskArtist(ctx, quizManager, int64(intNumberOfQuestions), bot, event)
							} else {
								title := "選択肢から選んでください。"
								replyMessage.AskTheNumberOfQuestionsMessage(ctx, bot, event, title)
							}
						} else if quizManager.Status == "artist" {
							artist := message.Text
							handler.InputArtistAndStartQuiz(ctx, quizManager, artist, bot, event)
						} else if quizManager.Status == "started" {
							if message.Text == "next" {
								handler.NextQuiz(ctx, quizManager, bot, event)
							} else if message.Text == "retry" {
								handler.Retry(ctx, quizManager, bot, event)
							} else if message.Text == "unknown_question" {
								handler.UnknownQuestion(ctx, quizManager, bot, event)
							} else if message.Text == "next_lyrics" {
								handler.NextLyrics(ctx, quizManager, bot, event)
							} else {
								handler.CheckAnswer(ctx, quizManager, message.Text, bot, event)
							}
						}
					}
				}
			}
		}
	})

	return r
}
