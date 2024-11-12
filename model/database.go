package model

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectToDatabase() *sql.DB {

	cfg := mysql.Config{
		User:   "admin",
		Passwd: "admin",
		Net:    "tcp",
		Addr:   "localhost",
		DBName: "inventario",
	}
	cfg.AllowNativePasswords = true

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to database")
	return db
}
