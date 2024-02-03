package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InitDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := "localhost" // または Docker Compose で定義されたサービス名
	dbPort := "3306"      // デフォルトのMySQLポート
	dbName := os.Getenv("MYSQL_DATABASE")

	// データベースへの接続文字列を構築
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	// データベースに接続
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// データベース接続をテスト
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database.")
	return DB
}
