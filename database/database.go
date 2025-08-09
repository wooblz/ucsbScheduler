package database

import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/wooblz/ucsbScheduler/models"
)

func CreateTable() error  {
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
            section_id INT REFERENCE
        )
    `
}
func InsertAllClasses(classes []models.Class) error  {

}

func ResetDB() error {

}

func 
