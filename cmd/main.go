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
