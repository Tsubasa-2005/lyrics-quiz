package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// ConnectDB SQLiteデータベースに接続または新規作成する関数
func ConnectDB(ctx context.Context) *sql.DB {
	// データベースファイルを作成または接続
	db, err := sql.Open("sqlite3", "data.sqlite3")
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Ping the database to verify connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return db
}

// ExecuteSQLFile SQLファイルを実行する関数
func ExecuteSQLFile(db *sql.DB, filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read SQL file: %v", err)
	}

	sqlStmt := string(content)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("failed to execute SQL file: %v", err)
	}
}

// migrateCmd データベースのテーブルを作成するcobraコマンド
func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 環境変数の設定
			os.Setenv("DB_PATH", "./example.db")

			// データベースに接続または新規作成
			ctx := context.Background()
			db := ConnectDB(ctx)
			defer db.Close()

			// SQLファイルを実行してテーブルを作成
			ExecuteSQLFile(db, "db/core.sql")

			fmt.Println("Database connected and tables created successfully!")
			return nil
		},
	}
	return cmd
}
