package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

func init() {
	dbUser := os.Getenv("SAUNA_USERNAME")
	dbPassword := os.Getenv("SAUNA_PW")
	dbDatabase := os.Getenv("DATABASE")
	dbHost := os.Getenv("MYSQL_HOST")
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbDatabase)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		panic(err) // 起動時に接続できなければ落とす
	}

	// 接続確認
	if err := db.Ping(); err != nil {
		panic(err)
	}

	Conn = db
	fmt.Println("MySQL Connect Success!!!")
}

func NewMySQLFromEnv() (*sql.DB, error) {
	dbUser := os.Getenv("SAUNA_USERNAME")
	dbPassword := os.Getenv("SAUNA_PW")
	dbDatabase := os.Getenv("DATABASE")
	dbHost := os.Getenv("MYSQL_HOST")
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbDatabase)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
