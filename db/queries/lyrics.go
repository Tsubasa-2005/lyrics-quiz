package queries

import (
	dbModel "lyrics-quiz/db"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CreateLyrics(c *gin.Context, lyrics dbModel.Lyrics) error {
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Create(&lyrics).Error; err != nil {
		return err
	}
	return nil
}

func GetLyrics(c *gin.Context, lyricsID string) (*dbModel.Lyrics, error) {
	db := c.MustGet("db").(*gorm.DB)

	var lyrics dbModel.Lyrics
	if err := db.Where("id = ?", lyricsID).First(&lyrics).Error; err != nil {
		return nil, err
	}
	return &lyrics, nil
}

func DeleteLyrics(c *gin.Context, lyricsID string) error {
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id = ?", lyricsID).Unscoped().Delete(&dbModel.Lyrics{}).Error; err != nil {
		return err
	}
	return nil
}
