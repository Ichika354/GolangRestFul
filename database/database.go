package database

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    var err error
    connStr := "user=postgres.sxntyesqqropbxdgipqe password=UpayeW4Ok8opNViy host=aws-0-ap-southeast-1.pooler.supabase.com port=5432 dbname=postgres"
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Database connected")
}

// dbConnectionString := fmt.Sprintf("user=postgres.juqrvjcwczdxuyzhuqoo password=hvPFh7fEdEKXd4QN host=aws-0-ap-southeast-1.pooler.supabase.com port=5432 dbname=postgres")