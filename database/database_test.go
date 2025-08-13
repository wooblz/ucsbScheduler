package database

import (
    "testing"
    "os"
    "log"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/wooblz/ucsbScheduler/models"
    "github.com/joho/godotenv"
)

func TestMain(m *testing.M)  {
    err := CreateTable()
    if err != nil  {
        log.Println(err)
    }

    code, err := run(m)
    if err != nil  {
        log.Println(err)
    }
    os.Exit(code)
}

func run(m *testing.M) (code int, err error)  {
    db_url := os.Getenv("TEST_DB_URL")
    db, err := sql.Open("postgres", db_url)
    if err != nil  {
        return -1, fmt.Errorf("Failed to connect to db: %v", err)
    }

    defer func()  {
        emptyTable, err := db.Exec("TRUNCATE TABLE test")
        if emptyTableErr != nil {
            return -1, fmt.Errorf("Failed to empty table: %v", err)
        }
        db.Close()
    }()
    return m.Run(), nil
}
