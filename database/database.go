package database

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
    var err error

    // Load environment variables from .env file
    err = godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get environment variables
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    dbname := os.Getenv("DB_NAME")

    // Build connection string
    connStr := "user=" + user + " password=" + password + " host=" + host + " port=" + port + " dbname=" + dbname

    // Open connection to database
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    // Ping database to ensure connection is established
    err = DB.Ping()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Database connected")
}
