package database

import (
    "testing"
    "fmt"
    "os"
    "log"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/wooblz/ucsbScheduler/models"
    //"github.com/joho/godotenv"
)

var db *sql.DB
func TestMain(m *testing.M)  {
    code, err := run(m)
    if err != nil  {
        log.Println(err)
    }
    os.Exit(code)
}

func run(m *testing.M) (code int, err error)  {
    db, err = StartDB("TEST_DB_URL")
    if err != nil  {
        return -1, fmt.Errorf("Failed to connect to db: %v", err)
    }
    _, _ = db.Exec("TRUNCATE TABLE classes, sections, time_locations RESTART IDENTITY CASCADE;")
    err = CreateTable(db)
    if err != nil  {
        log.Println(err)
    }
    defer func()  {
        _, _ = db.Exec("TRUNCATE TABLE classes, sections, time_locations RESTART IDENTITY CASCADE;")
        db.Close()
    }()
    return m.Run(), nil

}
func TestInsertClass(t *testing.T)  {
    _, _ = db.Exec("TRUNCATE TABLE classes, sections, time_locations RESTART IDENTITY CASCADE;")
    err := InsertAllClasses(Test1, db)
    if err != nil  {
        t.Fatalf("Failed to insert classes: %v", err)
    }

    queryClass,err := db.Prepare(`SELECT course_id, title, subject_area FROM classes WHERE course_id = $1`)
    if err != nil  {
        t.Fatalf("Failed prepare: %v", err)
    }
    defer queryClass.Close()

    querySection,err := db.Prepare(`SELECT id FROM sections WHERE course_id = $1`)
    if err != nil  {
        t.Fatalf("Failed prepare: %v", err)
    }
    defer querySection.Close()

    queryTime, err := db.Prepare(`SELECT room, building, days, begin_time, end_time FROM time_locations WHERE section_id = $1`)
    if err != nil  {
        t.Fatalf("Failed prepare: %v", err)
    }
    defer queryTime.Close()

    var count int
    err = db.QueryRow(`SELECT COUNT(*) FROM classes`).Scan(&count)
    if err != nil  {
        t.Fatalf("Failed to query count: %v", err)
    }
    if(count != 4)  {
        t.Fatalf("Miscount: %d", count)
    }

    var title, courseID, subjectArea string
    for _, class := range Test1  {
        mistake := false
        err = queryClass.QueryRow(class.CourseID).Scan(&courseID, &title, &subjectArea)
        if err != nil  {
            t.Fatalf("Failed to query row: %v",err)
        }
        if(courseID != class.CourseID || title != class.Title || subjectArea != class.SubjectArea)  {
            mistake = true 
        }
        c := models.Class{
            CourseID:    courseID,
            Title:       title,
            SubjectArea: subjectArea,
        }
        var sections []models.Section
        for _, s := range class.ClassSections  {
            sections = append(sections,s)
        }
        var id int
        err = querySection.QueryRow(class.CourseID).Scan(&id)
        if err != nil  {
            t.Fatalf("Failed to query row: %v", err)
        }
        row, err := queryTime.Query(id)
        if err != nil  { 
            t.Fatalf("Failed to query row: %v",err)
        }
        defer row.Close()
        for row.Next()  {
            found := false; 
            var room, building, days, begin_time, end_time string
            err := row.Scan(&room,&building,&days,&begin_time,&end_time)
            if err != nil  {
                t.Fatalf("Failed to query row: %v",err)
            }
            for i:=0; i < len(sections); i++  {
                if(room == sections[i].TimeLocations[0].Room && building == sections[i].TimeLocations[0].Building && 
                days == sections[i].TimeLocations[0].Days && begin_time == sections[i].TimeLocations[0].BeginTime && 
                end_time == sections[i].TimeLocations[0].EndTime)  {
                    sections = append(sections[:i],sections[i+1:]...)
                    found = true    
                    break
                }
            }
            s := models.Section {
                TimeLocations: []models.TimeLocation  {
                    {Room: room, Building: building, Days: days, BeginTime: begin_time, EndTime: end_time},
                },
            }
            c.ClassSections = append(c.ClassSections, s)
            if(!found)  {
                mistake = true
            }
        }
        if mistake  {
            t.Fatalf("Expected: \n%+v, Got: \n%+v",class,c)
        }
    }
    _, _ = db.Exec("TRUNCATE TABLE classes, sections, time_locations RESTART IDENTITY CASCADE;")
}

func TestQuery(t *testing.T)  {
    _, _ = db.Exec("TRUNCATE TABLE classes, sections, time_locations RESTART IDENTITY CASCADE;")
    err := InsertAllClasses(Test1, db)
    if err != nil  {
        t.Fatalf("Failed to Insert Classes: %v", err)
    }
    classes, err := QueryTitle("CMPSC", db)
    if err != nil  {
        t.Fatalf("Failed to Query Classes: \n%v", err)
    }
    fmt.Printf("%+v", classes)
    _, _ = db.Exec("TRUNCATE TABLE classes, sections, time_locations RESTART IDENTITY CASCADE;")
}
