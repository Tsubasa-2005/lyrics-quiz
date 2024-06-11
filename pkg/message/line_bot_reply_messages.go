package message

import (
	"fmt"
	"lyrics-quiz/pkg/infra/rdb"

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

func ErrorCreatingProblems(c *gin.Context, bot *linebot.Client, event *linebot.Event) {
	reply := "問題作成中にエラーが起きました。もう一度アーティストの名前を入力してください。(何度も続く場合アーティストの名前を変更してみてください。)"
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

func ProblemStatement(c *gin.Context, bot *linebot.Client, event *linebot.Event, quizManager rdb.QuizManager, repo *rdb.Queries) {
	nowQuestion := quizManager.QuizCount
	nowLyrics := quizManager.LyricsCount
	lyrics, err := repo.GetLyrics(c, rdb.GetLyricsParams{
		QuizManagerID:  quizManager.UserID,
		QuestionNumber: nowQuestion,
		Count:          nowLyrics,
	})
	if err != nil {
		Error(c, bot, event)
		return
	}
	choices, err := repo.GetChoices(c, rdb.GetChoicesParams{
		QuizManagerID:  quizManager.UserID,
		QuestionNumber: nowQuestion,
	})
	question := lyrics.Lyrics
	choice1 := choices.Choice1
	choice2 := choices.Choice2
	choice3 := choices.Choice3
	choice4 := choices.Choice4

	title := fmt.Sprintf("第%d問 歌詞の曲名を選択してください。", nowQuestion)
	contents := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   title,
					Size:   linebot.FlexTextSizeTypeMd,
					Weight: linebot.FlexTextWeightTypeBold,
					Wrap:   true,
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: question,
					Size: linebot.FlexTextSizeTypeMd,
					Wrap: true,
				},
				&linebot.ButtonComponent{
					Action: linebot.NewMessageAction(choice1, choice1),
				},
				&linebot.ButtonComponent{
					Action: linebot.NewMessageAction(choice2, choice2),
				},
				&linebot.ButtonComponent{
					Action: linebot.NewMessageAction(choice3, choice3),
				},
				&linebot.ButtonComponent{
					Action: linebot.NewMessageAction(choice4, choice4),
				},
				&linebot.ButtonComponent{
					Action: linebot.NewMessageAction("わからない", "unknown_question"),
				},
				&linebot.ButtonComponent{
					Action: linebot.NewMessageAction("歌詞の続きを見る。", "next_lyrics"),
				},
			},
		},
	}

	message := linebot.NewFlexMessage("問題と選択肢", contents)

	if _, err := bot.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func CorrectAnswer(c *gin.Context, bot *linebot.Client, event *linebot.Event) {
	reply := "正解です！"
	template := linebot.NewButtonsTemplate(
		"",    // thumbnailImageUrl
		reply, // title
		"次の操作を選んでください。", // text
		linebot.NewMessageAction("次の問題を解く", "next"),
		linebot.NewMessageAction("クイズを終了する", "e"),
	)

	message := linebot.NewTemplateMessage("次の問題", template)

	if _, err := bot.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func IncorrectAnswer(c *gin.Context, bot *linebot.Client, event *linebot.Event) {
	reply := "不正解です。"
	template := linebot.NewButtonsTemplate(
		"",             // thumbnailImageUrl
		reply,          // title
		"次の操作を選んで下さい。", // text
		linebot.NewMessageAction("次の問題を解く", "next"),
		linebot.NewMessageAction("クイズを終了する", "e"),
		linebot.NewMessageAction("もう一度挑戦する", "retry"),
		linebot.NewMessageAction("歌詞の続きを見る", "next_lyrics"),
	)

	message := linebot.NewTemplateMessage("不正解", template)

	if _, err := bot.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func EndQuiz(c *gin.Context, bot *linebot.Client, event *linebot.Event) {
	reply := "正解!! これで出題を終わります。お疲れ様でした！"
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}

func Unknown(c *gin.Context, bot *linebot.Client, event *linebot.Event, answer string) {
	reply := fmt.Sprintf("正解は、%sです。", answer)
	template := linebot.NewButtonsTemplate(
		"",           // thumbnailImageUrl
		reply,        // title
		"操作を選んでください", // text
		linebot.NewMessageAction("次の問題を解く", "next"),
		linebot.NewMessageAction("クイズを終了する", "e"),
	)

	message := linebot.NewTemplateMessage("次の問題", template)

	if _, err := bot.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
		c.Error(errors.Wrap(err))
	}
}
