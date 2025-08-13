package database

import (
    "fmt"
    "log"
    "os"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/wooblz/ucsbScheduler/models"
    "github.com/joho/godotenv"
)

func CreateTable() error  {
    err := godotenv.Load()
    if err != nil   {
        return err
    }
    db_url := os.Getenv("DB_URL")
    db, err := sql.Open("postgres", db_url)
    defer db.Close()
    if err != nil  {
        log.Printf("Unable to open server: %v", err)
        return err 
    }  else  {
        log.Println("Server Opened")
    }
    connectivity := db.Ping()
    if connectivity != nil  {
        log.Printf("Unable to ping server: %v", err)
    }  else  {
        log.Println("Server pinged")
    }
    TableCreate := `
        CREATE TABLE classes (
            course_id TEXT PRIMARY KEY, 
            title TEXT NOT NULL,
            subject_area TEXT NOT NULL,
        );
        CREATE TABLE sections (
            id SERIAL PRIMARY KEY,
            course_id TEXT REFERENCES classes(course_id) ON DELETE CASCADE
        );
        CREATE TABLE time_locations (
            id SERIAL PRIMARY KEY,
            section_id INT REFERENCE sections(id) ON DELETE CASCADE,
            room TEXT,
            building TEXT,
            days TEXT,
            begin_time TIME,
            end_time TIME
        );
    `
    _, err = db.Exec(TableCreate)
    if err != nil  {
        log.Printf("Failed to create table: %v", err)
        return err
    }
    return nil
}
func InsertAllClasses(classes []models.Class) error  {

}

func ResetDB() error {

}

func 
