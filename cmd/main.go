package main

import (
	"context"
	"log"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	if err := rootCmd(ctx).Execute(); err != nil {
		log.Fatal(err)
	}
}

func rootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "cli",
	}

	cmd.AddCommand(
		serverCmd(ctx),
		migrateCmd(),
	)

	return cmd
}

/*func main() {
	db, err := gorm.Open("sqlite3", "data.sqlite3")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
		return
	}
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Migrate the schema
	if err := db.AutoMigrate(&dbModel.QuizCounter{}).Error; err != nil {
		log.Fatalf("Failed to migrate QuizCounter model: %v", err)
		return
	}
	if err := db.AutoMigrate(&dbModel.LyricsCounter{}).Error; err != nil {
		log.Fatalf("Failed to migrate LyricsCounter model: %v", err)
		return
	}
	if err := db.AutoMigrate(&dbModel.Lyrics{}).Error; err != nil {
		log.Fatalf("Failed to migrate Ltrics model: %v", err)
		return
	}
	if err := db.AutoMigrate(&dbModel.Answer{}).Error; err != nil {
		log.Fatalf("Failed to migrate Answer model: %v", err)
		return
	}
}*/
