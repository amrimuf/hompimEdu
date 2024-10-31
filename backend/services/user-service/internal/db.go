package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

func InitDB() (*sql.DB, error) {
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")

    // More explicit connection string format
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    var db *sql.DB
    var err error

    log.Printf("Attempting to connect to database with DSN: host=%s port=%s user=%s dbname=%s",
        dbHost, dbPort, dbUser, dbName) // Log connection attempt (omit password)

    // Attempt to connect to the database with retries
    for retries := 5; retries > 0; retries-- {
        db, err = sql.Open("pgx", connStr)
        if err == nil {
            // Check if we can actually ping the database
            err = db.Ping()
            if err == nil {
                log.Println("Connected to the database successfully")
                return db, nil
            }
        }

        log.Printf("Failed to connect to database (retries left: %d): %v", retries-1, err)
        time.Sleep(2 * time.Second)
    }

    return nil, fmt.Errorf("failed to connect to database after 5 attempts: %v", err)
}
