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
    err = ResetDB(db)
    if err != nil  {
        return -1, fmt.Errorf("Failed to  reset db: %v", err)
    }
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
    err := ResetDB(db)
    if err != nil  {
        t.Fatalf("Failed to reset: %v", err)
    }
    err = InsertAllClasses(Test1, db)
    if err != nil  {
        t.Fatalf("Failed to insert classes: %v", err)
    }

    queryClass,err := db.Prepare(`SELECT course_id, title, subject_area, room, building, days, begin_time, end_time FROM classes WHERE course_id = $1`)
    if err != nil  {
        t.Fatalf("Failed prepare: %v", err)
    }
    defer queryClass.Close()

    querySection,err := db.Prepare(`SELECT id, enroll_code FROM sections WHERE course_id = $1`)
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

    var title, courseID, subjectArea, room, building, days, begin_time, end_time string
    for _, class := range Test1  {
        mistake := false
        err = queryClass.QueryRow(class.CourseID).Scan(&courseID, &title, &subjectArea, &room, &building, &days, &begin_time, &end_time)
        if err != nil  {
            t.Fatalf("Failed to query row: %v",err)
        }
        if(courseID != class.CourseID || title != class.Title || subjectArea != class.SubjectArea || room != class.Room || building != class.Building || days != class.Days || begin_time != class.BeginTime || end_time != class.EndTime)  {
            mistake = true 
        }
        c := models.Class{
            CourseID:    courseID,
            Title:       title,
            SubjectArea: subjectArea,
            Room: room,
            Building: building, 
            Days: days, 
            BeginTime: begin_time, 
            EndTime: end_time,
        }
        var sections []models.Section
        for _, s := range class.ClassSections  {
            sections = append(sections,s)
        }
        if len(class.ClassSections) == 0 {
            continue
        }
        var id int
        var enrollCode string
        err = querySection.QueryRow(class.CourseID).Scan(&id, &enrollCode)
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
                EnrollCode : enrollCode, 
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
    err = ResetDB(db)
    if err != nil  {
        t.Fatalf("Failed to reset: %v", err)
    }
}

func TestQuery(t *testing.T)  {
    err := ResetDB(db)
    if err != nil  {
        t.Fatalf("Failed to reset: %v", err)
    }
    err = InsertAllClasses(Test1, db)
    if err != nil  {
        t.Fatalf("Failed to Insert Classes: %v", err)
    }
    classes, err := QueryTitle("CMPSC", db)
    if err != nil  {
        t.Fatalf("Failed to Query Classes: \n%v", err)
    }
    for _, c := range classes {
		fmt.Printf("Class: %s, Title: %s, Enroll: %s, Time: %s %s-%s\n", c.CourseID, c.Title, c.EnrollCode, c.Days, c.BeginTime, c.EndTime)
		for _, s := range c.ClassSections {
			fmt.Printf("  - Section Enroll: %s\n", s.EnrollCode)
			for _, tl := range s.TimeLocations {
				fmt.Printf("    -> Location: %s-%s, Time: %s %s-%s\n", tl.Building, tl.Room, tl.Days, tl.BeginTime, tl.EndTime)
			}
		}
		fmt.Println("--------------------")
	}
    err = ResetDB(db)
    if err != nil  {
        t.Fatalf("Failed to reset: %v", err)
    }
}
