package queries

import (
	dbModel "lyrics-quiz/db"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UpdateQuizCount(c *gin.Context, quizCount dbModel.QuizCounter) error {
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Model(&dbModel.QuizCounter{}).Where("id = ?", quizCount.ID).Updates(dbModel.QuizCounter{
		Count: quizCount.Count,
	}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateLyricsCount(c *gin.Context, lyricsCount dbModel.LyricsCounter) error {
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Model(&dbModel.LyricsCounter{}).Where("id = ?", lyricsCount.ID).Updates(dbModel.LyricsCounter{
		Count: lyricsCount.Count,
	}).Error; err != nil {
		return err
	}
	return nil
}
