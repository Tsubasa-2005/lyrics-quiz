package main

import (
	"fmt"
	dbModel "lyrics-quiz/db"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate and initialize",
		Short: "Migrate the database and initialize database",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := gorm.Open("sqlite3", "data.sqlite3")
			if err != nil {
				return fmt.Errorf("failed to open database: %v", err)
			}
			defer func(db *gorm.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			// Migrate the schema
			if err := db.AutoMigrate(&dbModel.QuizCounter{}).Error; err != nil {
				return fmt.Errorf("failed to migrate QuizCounter model: %v", err)
			}
			if err := db.AutoMigrate(&dbModel.LyricsCounter{}).Error; err != nil {
				return fmt.Errorf("failed to migrate LyricsCounter model: %v", err)
			}
			if err := db.AutoMigrate(&dbModel.Lyrics{}).Error; err != nil {
				return fmt.Errorf("failed to migrate Lyrics model: %v", err)
			}
			if err := db.AutoMigrate(&dbModel.Answer{}).Error; err != nil {
				return fmt.Errorf("failed to migrate Answer model: %v", err)
			}

			err = initializeDataBase(db)

			return err
		},
	}
	return cmd
}

func initializeDataBase(db *gorm.DB) error {
	quizCount := dbModel.QuizCounter{
		QuizCounterID: 1, Count: 1,
	}
	var err error

	if err = db.Create(&quizCount).Error; err != nil {
		return fmt.Errorf("failed to initialize quizCounter model: %v", err)
	}

	lyricsCount := dbModel.LyricsCounter{
		LyricsCounterID: 1, Count: 1,
	}

	if err = db.Create(&lyricsCount).Error; err != nil {
		return fmt.Errorf("failed to initialize lyricsCounter model: %v", err)
	}
	return nil
}
