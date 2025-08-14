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
    err := godotenv.Load("../.env")
    if err != nil   {
        return err
    }

    db_url := os.Getenv("DB_URL")
    db, err := sql.Open("postgres", db_url)
    if err != nil  {
        return err 
    }
    defer db.Close()
    connectivity := db.Ping()
    if connectivity != nil  {
        return fmt.Errorf("Unable to ping server: %v", err)
    }  
    TableCreate := `
        CREATE TABLE classes (
            course_id TEXT PRIMARY KEY, 
            title TEXT NOT NULL,
            subject_area TEXT NOT NULL
            tsv tsvector
        );
        CREATE TABLE sections (
            id SERIAL PRIMARY KEY,
            course_id TEXT REFERENCES classes(course_id) ON DELETE CASCADE
        );
        CREATE TABLE time_locations (
            id SERIAL PRIMARY KEY,
            section_id INT REFERENCES sections(id) ON DELETE CASCADE,
            room TEXT,
            building TEXT,
            days TEXT,
            begin_time TIME,
            end_time TIME
        );
    `
    _, err = db.Exec(TableCreate)
    if err != nil  {
        return err
    }
    return nil
}
func InsertAllClasses(classes []models.Class) error  {
    insert_class,err := db.Prepare("INSERT INTO classes VALUES ($1, $2, $3)")
    if err != nil {
        return err
    }
    defer insert_class.Close()
    insert_section,err  := db.Prepare("INSERT INTO section (course_id) VALUES ($1) RETURNING id")
    if err != nil  {
        return err
    }
    defer insert_section.Close()
    insert_time, err := db.Prepare("INSERT INTO time_locations 
        (section_id, room, building, days, begin_time, end_time) 
        VALUES ($1, $2, $3, $4, $5, $6)"
    )
    if  err != nil  {
        return err
    }
    defer insert_time.Close()
    for _, v := range classes {
        _, err = insert_class.Exec(v.CourseID, v.Title, v.SubjectArea)
        if  err != nil  {
            return err 
        }
        for _, w := range v.ClassSections  {
            var SectionID int
            _, err = insert_section.QueryRow(v.CourseID).Scan(&SectionID)
            if err != nil  {
                return err
            }
            for _, x := range w.TimeLocataions  {
                _, err = insert_time.Exec(SectionID, x.Room,x.Building, x.Days, x.BeginTime, x.EndTime)
                if err != nil  {
                    return err
                }
            }
        }
    }
    db.Exec("UPDATE classes SET tsv = to_tsvector('english', coalesce(title,''))")
    return nil
}

func ResetDB() error {
    _, err := db.Exec("TRUNCATE classes RESTART IDENTITY CASCADE")
    if err != nil  {
        return err
    }
    return nil
}

func QueryTitle(Query string) (classes []Class, error) {

}
