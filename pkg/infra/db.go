package infra

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattn/go-sqlite3"
)

// ConnectDB SQLiteデータベースに接続または新規作成する関数
func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "data.sqlite3")
	if err != nil {
		panic(err)
	}

	return db
}

// CheckConnectDB データベース接続を確認する関数
func CheckConnectDB() {
	dbConn := ConnectDB()
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		panic(fmt.Sprintf("failed to ping database: %v", err))
	}

	fmt.Println("Database connected successfully!")
}
