package queries

import (
	dbModel "lyrics-quiz/db"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CreateAnswer(c *gin.Context, answer dbModel.Answer) error {
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Create(&answer).Error; err != nil {
		return err
	}
	return nil
}

func GetAnswer(c *gin.Context, ID string) (*dbModel.Answer, error) {
	db := c.MustGet("db").(*gorm.DB)

	var answer dbModel.Answer
	if err := db.Where("id = ?", ID).First(&answer).Error; err != nil {
		return nil, err
	}
	return &answer, nil
}

func DeleteAnswer(c *gin.Context, ID string) error {
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id = ?", ID).Unscoped().Delete(&dbModel.Answer{}).Error; err != nil {
		return err
	}

	return nil
}
