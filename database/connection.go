package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	// Define the MySQL DSN (adjust user/password as per your Docker configuration)
	dsn := "rani:1234@tcp(localhost:3306)/user_management"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Verify the connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Database connected successfully")
}
