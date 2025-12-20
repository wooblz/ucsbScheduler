package database

import (
    "fmt"
    "os"
    "encoding/json"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/wooblz/ucsbScheduler/models"
)

const url string = "DB_URL"

func StartDB(name string) (*sql.DB, error)  {
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
            course_id INT REFERENCES classes(id) ON DELETE CASCADE
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
}
func GetSelectedClass(enrollCode string, db *sql.DB) (models.Class, error) {
	var c models.Class
	var sEnroll, sRoom, sBuild, sDays, sBegin, sEnd string

	row := db.QueryRow(`
		SELECT 
			c.course_id, c.title, c.subject_area, c.enroll_code, c.room, c.building, c.days, c.begin_time, c.end_time,
			s.enroll_code, tl.room, tl.building, tl.days, tl.begin_time, tl.end_time
		FROM sections s
		JOIN classes c ON s.course_id = c.id
		JOIN time_locations tl ON tl.section_id = s.id
		WHERE s.enroll_code = $1
	`, enrollCode)

	err := row.Scan(
		&c.CourseID, &c.Title, &c.SubjectArea, &c.EnrollCode, &c.Room, &c.Building, &c.Days, &c.BeginTime, &c.EndTime,
		&sEnroll, &sRoom, &sBuild, &sDays, &sBegin, &sEnd,
	)

	if err == nil {
		c.ClassSections = []models.Section{
			{
				EnrollCode: sEnroll,
				TimeLocations: []models.TimeLocation{
					{Room: sRoom, Building: sBuild, Days: sDays, BeginTime: sBegin, EndTime: sEnd},
				},
			},
		}
		return c, nil
	}

	row = db.QueryRow(`
		SELECT course_id, title, subject_area, enroll_code, room, building, days, begin_time, end_time
		FROM classes 
		WHERE enroll_code = $1
	`, enrollCode)

	err = row.Scan(
		&c.CourseID, &c.Title, &c.SubjectArea, &c.EnrollCode, &c.Room, &c.Building, &c.Days, &c.BeginTime, &c.EndTime,
	)

	if err != nil {
		return models.Class{}, fmt.Errorf("enroll code %s not found", enrollCode)
	}

	c.ClassSections = []models.Section{}
	return c, nil
}
