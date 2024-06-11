package message

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/taxio/errors"
)

func AskQuizTypeMessage(c *gin.Context, bot *linebot.Client, event *linebot.Event, title string) {
	template := linebot.NewButtonsTemplate(
		"",    // thumbnailImageUrl
		title, // title
		"出題するクイズのタイプを選択してください。", // text
		linebot.NewMessageAction("歌詞クイズ(normal)", "normal"),
		linebot.NewMessageAction("歌詞クイズ(hard)", "hard"),
	)

	message := linebot.NewTemplateMessage("クイズのタイプ", template)

	if _, err := bot.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func AskTheNumberOfQuestionsMessage(c *gin.Context, bot *linebot.Client, event *linebot.Event, title string) {
	template := linebot.NewButtonsTemplate(
		"",    // thumbnailImageUrl
		title, // title
		"出題する問題数を選択してください。", // text
		linebot.NewMessageAction("1問", "1"),
		linebot.NewMessageAction("5問", "5"),
		linebot.NewMessageAction("10問", "10"),
	)

	message := linebot.NewTemplateMessage("問題数", template)

	if _, err := bot.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func AskArtist(c *gin.Context, bot *linebot.Client, event *linebot.Event) {
	reply := "問題を出してほしいアーティストの名前を入力してください。(アーティストの名前を入力すると問題を生成します。この処理には、少し時間がかかります。)"
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func Error(c *gin.Context, bot *linebot.Client, event *linebot.Event) {
	reply := "エラーが発生しました。もう一度やり直してください。"
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func FinishedInitialize(c *gin.Context, bot *linebot.Client, event *linebot.Event) {
	reply := "初期設定が完了しました。"
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func FinishedQuiz(c *gin.Context, bot *linebot.Client, event *linebot.Event) {
	reply := "クイズを終了しました。"
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func Parroting(c *gin.Context, bot *linebot.Client, event *linebot.Event, message *linebot.TextMessage) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}
