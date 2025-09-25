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
    /*err = ResetDB(db)
    if err != nil  {
        return -1, fmt.Errorf("Failed to  reset db: %v", err)
    }*/
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
         package database

import (
    "fmt"
    "os"
    "encoding/json"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/wooblz/ucsbScheduler/models"
    "github.com/joho/godotenv"
)

const url string = "DB_URL"

func StartDB(name string) (*sql.DB, error)  {
    err := godotenv.Load("../.env")
    if err != nil   {
        return nil, err
    }
    db_url := os.Getenv(name)
    db, err := sql.Open("postgres", db_url)
    if err != nil  {
        return nil, err
    }
    connectivity := db.Ping()
    if connectivity != nil  {
        return nil, fmt.Errorf("Unable to ping server: %v", connectivity)
    }
    return db, nil
}
func CreateTable(db *sql.DB) error  {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS classes (
            id SERIAL PRIMARY KEY,
            course_id TEXT,
            title TEXT NOT NULL,
            subject_area TEXT NOT NULL,
            enroll_code TEXT,
            room TEXT,
            building TEXT,
            days TEXT,
            begin_time TEXT,
            end_time TEXT,
            tsv tsvector
        )`)
    if err != nil  {
        return err
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS sections (
            id SERIAL PRIMARY KEY,
            enroll_code TEXT,
            course_id TEXT REFERENCES classes(id) ON DELETE CASCADE
        )`)
    if err != nil  {
        return err
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS time_locations (
            id SERIAL PRIMARY KEY,
            section_id INT REFERENCES sections(id) ON DELETE CASCADE,
            room TEXT,
            building TEXT,
            days TEXT,
            begin_time TEXT,
            end_time TEXT
        )`)
    if err != nil  {
        return err
    }
    return nil
}
func InsertAllClasses(classes []models.Class, db *sql.DB) error  {
    insert_class,err := db.Prepare("INSERT INTO classes (course_id, title, subject_area, enroll_code, room, building, days, begin_time, end_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id")
    if err != nil {
        return err
    }
    defer insert_class.Close()
    insert_section,err  := db.Prepare("INSERT INTO sections (course_id, enroll_code) VALUES ($1, $2) RETURNING id")
    if err != nil  {
        return err
    }
    defer insert_section.Close()
    insert_time, err := db.Prepare(`INSERT INTO time_locations
        (section_id, room, building, days, begin_time, end_time)
        VALUES ($1, $2, $3, $4, $5, $6)`)
    if  err != nil  {
        return err
    }
    defer insert_time.Close()
    for _, v := range classes {
        var ClassID int
        err = insert_class.QueryRow(v.CourseID, v.Title, v.SubjectArea, v.EnrollCode, v.Room, v.Building, v.Days, v.BeginTime, v.EndTime).Scan(&ClassID)
        if  err != nil  {
            return err
        }
        for _, w := range v.ClassSections  {
            var SectionID int
            err = insert_section.QueryRow(ClassID,w.EnrollCode).Scan(&SectionID)
            if err != nil  {
                return err
            }
            for _, x := range w.TimeLocations  {
                _, err = insert_time.Exec(SectionID, x.Room,x.Building, x.Days, x.BeginTime, x.EndTime)
                if err != nil  {
                    return err
                }
            }
        }
    }
    db.Exec(`UPDATE classes SET tsv =
        setweight(to_tsvector('english', coalesce(course_id, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(title, '')), 'B')`)
    return nil
}

func ResetDB(db *sql.DB) error {
    _, err := db.Exec("TRUNCATE TABLE classes RESTART IDENTITY CASCADE")
    if err != nil  {
        return err
    }
    return nil
}

func QueryTitle(statement string, db *sql.DB) ([]models.Class, error) {
    query_line, err := db.Prepare(`
        SELECT
            c.course_id,
            c.title,
            c.subject_area,
            c.enroll_code,
            c.room,
            c.building,
            c.days,
            c.begin_time,
            c.end_time,
            (
                SELECT COALESCE(json_agg(s_obj), '[]'::json)
                FROM (
                    SELECT
                        json_build_object(
                            'enrollCode', s.enroll_code,
                            'timeLocations', (
                                SELECT COALESCE(json_agg(tl_obj), '[]'::json)
                                FROM (
                                    SELECT json_build_object(
                                        'room', tl.room,
                                        'building', tl.building,
                                        'days', tl.days,
                                        'beginTime', tl.begin_time,
                                        'endTime', tl.end_time
                                    ) as tl_obj
                                    FROM time_locations tl
                                    WHERE tl.section_id = s.id
                                    ORDER BY tl.id
                                ) tl_sub
                            )
                        ) as s_obj
                    FROM sections s
                    WHERE s.course_id = c.id
                    ORDER BY s.enroll_code
                ) s_sub
            ) as classSections
        FROM classes c
        WHERE c.tsv @@ plainto_tsquery('english', $1);
    `)
    if err != nil  {
        return nil, err
    }
    defer query_line.Close()
    rows, err := query_line.Query(statement)
    if err != nil {
        return nil, err
    }
    var classes []models.Class
    for rows.Next()  {
        var courseID, title, subjectArea, enrollCode, room, building, days, begin_time, end_time string
        var jsonData []byte
        err = rows.Scan(&courseID, &title, &subjectArea, &enrollCode, &room, &building, &days, &begin_time, &end_time, &jsonData)
        if err != nil  {
            return nil, err
        }
        var sections []models.Section
        err = json.Unmarshal(jsonData,&sections)
        if err != nil  {
            return nil, err
        }
        c := models.Class{
            CourseID: courseID,
            Title: title,
            SubjectArea: subjectArea,
            EnrollCode : enrollCode,
            Room: room,
            Building: building,
            Days: days,
            BeginTime: begin_time,
            EndTime: end_time,
            ClassSections: sections,
        }
        classes = append(classes,c)
    }
    return classes, nil
}   t.Fatalf("Failed to query row: %v", err)
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
