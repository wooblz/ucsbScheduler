package database

import (
    "testing"
    "fmt"
    "os"
    "log"
    //"database/sql"
    _ "github.com/lib/pq"
    //"github.com/wooblz/ucsbScheduler/models"
    //"github.com/joho/godotenv"
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
    
    db, err := StartDB("TEST_DB_URL")
    if err != nil  {
        return -1, fmt.Errorf("Failed to connect to db: %v", err)
    }

    defer func()  {
        _, _ = db.Exec("TRUNCATE TABLE test")
        db.Close()
    }()
    return m.Run(), nil

}

