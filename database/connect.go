package database

import (
    "database/sql"
    "fmt"
    "log"

    
    _ "github.com/lib/pq"
)

var DB *sql.DB

func init() {

    dbConnectionString := fmt.Sprintf("user=postgres.juqrvjcwczdxuyzhuqoo password=hvPFh7fEdEKXd4QN host=aws-0-ap-southeast-1.pooler.supabase.com port=5432 dbname=postgres")

    var err error
    DB, err = sql.Open("postgres", dbConnectionString)
    if err != nil {
        log.Fatal(err)
    }
}
