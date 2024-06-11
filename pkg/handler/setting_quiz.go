package handler

import (
	"fmt"
	"lyrics-quiz/pkg/infra/rdb"
	"lyrics-quiz/pkg/lyrics"
	"lyrics-quiz/pkg/message"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func AskQuizType(c *gin.Context, quizManager rdb.QuizManager, bot *linebot.Client, event *linebot.Event) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	quizManager.Status = "quiz_type"
	err := repo.UpdateQuizManager(c, rdb.UpdateQuizManagerParams{
		TheNumberOfQuestions: quizManager.TheNumberOfQuestions,
		QuizCount:            quizManager.QuizCount,
		LyricsCount:          quizManager.LyricsCount,
		Status:               quizManager.Status,
		Type:                 quizManager.Type,
		UserID:               quizManager.UserID,
	})
	if err != nil {
		message.Error(c, bot, event)
	}
	title := "クイズ開始!"
	message.AskQuizTypeMessage(c, bot, event, title)
}

func InputQuizTypeAndAskTheNumberOfQuestions(c *gin.Context, quizManager rdb.QuizManager, quizType string, bot *linebot.Client, event *linebot.Event) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	quizManager.Type = quizType
	quizManager.Status = "the_number_of_questions"
	err := repo.UpdateQuizManager(c, rdb.UpdateQuizManagerParams{
		TheNumberOfQuestions: quizManager.TheNumberOfQuestions,
		QuizCount:            quizManager.QuizCount,
		LyricsCount:          quizManager.LyricsCount,
		Status:               quizManager.Status,
		Type:                 quizManager.Type,
		UserID:               quizManager.UserID,
	})
	if err != nil {
		message.Error(c, bot, event)
	}

	title := "問題数を選択してください。"
	message.AskTheNumberOfQuestionsMessage(c, bot, event, title)
}

func InputTheNumberOfQuestionsAndAskArtist(c *gin.Context, quizManager rdb.QuizManager, theNumberOfQuestions int64, bot *linebot.Client, event *linebot.Event) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	quizManager.TheNumberOfQuestions = theNumberOfQuestions
	quizManager.Status = "artist"
	err := repo.UpdateQuizManager(c, rdb.UpdateQuizManagerParams{
		TheNumberOfQuestions: quizManager.TheNumberOfQuestions,
		QuizCount:            quizManager.QuizCount,
		LyricsCount:          quizManager.LyricsCount,
		Status:               quizManager.Status,
		Type:                 quizManager.Type,
		UserID:               quizManager.UserID,
	})
	if err != nil {
		message.Error(c, bot, event)
	}

	message.AskArtist(c, bot, event)
}

func InputArtistAndStartQuiz(c *gin.Context, quizManager rdb.QuizManager, artist string, bot *linebot.Client, event *linebot.Event) {
	dbConn := c.MustGet("db").(rdb.DBTX)
	repo := rdb.New(dbConn)

	quizManager.Status = "started"
	err := repo.UpdateQuizManager(c, rdb.UpdateQuizManagerParams{
		TheNumberOfQuestions: quizManager.TheNumberOfQuestions,
		QuizCount:            quizManager.QuizCount,
		LyricsCount:          quizManager.LyricsCount,
		Status:               quizManager.Status,
		Type:                 quizManager.Type,
		UserID:               quizManager.UserID,
	})
	if err != nil {
		message.Error(c, bot, event)
	}

	err = repo.CreateArtist(c, rdb.CreateArtistParams{
		QuizManagerID: quizManager.UserID,
		Artist:        artist,
	})
	if err != nil {
		message.Error(c, bot, event)
	}

	trackNames, err := lyrics.GetArtistTopTracksBySpotifyAPI(c, repo, artist, quizManager, bot, event)
	if err != nil {
		message.Error(c, bot, event)
	}

	createdLyrics, err := lyrics.GetLyrics(c, repo, trackNames, quizManager)
	if err != nil {
		message.Error(c, bot, event)
	}
	fmt.Println("Lyrics Parts[0]:", createdLyrics[0])
}
