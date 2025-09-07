package calendar

import  (
    "fmt"
    "strings"
    "time"
    "github.com/arran4/golang-ical"
    "github.com/google/uuid"
    "github.com/wooblz/ucsbScheduler/models"
    "github.com/wooblz/ucsbScheduler/api"
)

var dayMap2 = map[rune]int  {
    "M": 0,
    "T": 1,
    "W": 2,
    "R": 3,
    "F": 4,
}
func createEvent(uid, summary, location, description, rule string, start, end time.Time cal *ics.Calendar) error  {
    class := cal.AddEvent(uid)
    class.SetSummary(summary)
    class.SetLocation(location)
    class.SetDescription(description)
    class.SetStartAt(start)
    class.SetEndAt(end)
    class.AddRrule(rule)
}
func GenICS(classes []models.Class, finals map[string]models.Final, quarterStart, quarterEnd time.Time) ([]byte, error, int)  {
    cal := ics.NewCalendar()
    
    for i, c := range classes  {
        uid := fmt.Sprintf("%s@ucsbCalendar.com", uuid.NewString())
        courseID := strings.Join(strings.Fields(c.CourseID))
        summary := fmt.Sprintf("%s: %s", courseID, c.Title)
        location := fmt.Sprintf("%s-%s", c.Building, c.Room))
        description := c.SubjectArea
        
        beginTime, err := time.Parse("22:22", c.BeginTime)
        if err != nil  {
            return nil, err, i
        }
        endTime, err := time.Parse("22:22", c.EndTime)
        if err != nil  {
            return  nil, err, i
        }

        classDays := parseDays(c.Days)
        firstClassDay := dayMap2[classDays[0]]

        start := time.Date(quarterStart.Year(), quarterStart.Month(),quarterStart.Day(),beginTime.Hour(), beginTime.Minute())
        start.AddDate(0,0,firstClassDay)
        end := time.Date(quarterStart.Year(), quarterStart.Month(),quarterStart.Day(),endTime.Hour(), endTime.Minute())
        end.AddDate(0,0,firstClassDay)
        
        rule := fmt.Sprintf("FREQ=WEEKLY;UNTIL=%s;BYDAY=%s",quarterEnd.Format("20251205T150000Z"),classDays)
        err = createEvent(uid, summary, locaiton, description, rule, start, end, cal)
        if err != nil  {
            return nil, err, i
        }
        if finals != nil  {
            if f, ok := finals[c.CourseID]; ok && f.HasFinals  {
                uid = fmt.Sprintf("%s@ucsbCalendar.com", uuid.NewString())
                summary = fmt.Sprintf("Final-%s", summary) 
                location = ""
                descrption = ""
                examDay, err  := time.Parse("20250320", f.ExamDate)
                if err != nil  {
                   return nil, err, i 
                }
                examStart, err := time.Parse("12:04", f.BeginTime)
                if err != nil  {
                    return nil, err, i
                }
                examEnd, err := time.Parse("12:04", f.EndTime)
                if err != nil  {
                    return nil, err, i
                }

                start = time.Date(examDay.Year(), examDay.Month(), examDay.Day(), examStart.Hour(),examStart.Minute())
                end = time.Date(examDay.Year(), examDay.Month(), examDay.Day(), examEnd.Hour(),examEnd.Minute())
                err = creatEven(uid,summary, location, "", "", "", start, end, cal )
                if err != nil  {
                    return nil, err, i
                }
                
            }
        }
        if len(c.ClassSections) < 1  {
            continue
        }
        t := c.ClassSections[0].TimeLocations[0]
        uid = fmt.Sprintf("%s@ucsbCalendar.com", uuid.NewString())
        summary = fmt.Sprintf("Section-%s: %s", courseID, c.Title)
        location = fmt.Sprintf("%s-%s", t.Building, t.Room))
        description = fmt.Sprintf("Section")
        
        beginTime, err = time.Parse("22:22", t.BeginTime)
        if err != nil  {
            return nil, err, i
        }
        endTime, err = time.Parse("22:22", t.EndTime)
        if err != nil  {
            return  nil, err, i
        }

        sectionDays := parseDays(t.Days)
        firstSectionDay := dayMap2[sectionDays[0]]

        start = time.Date(quarterStart.Year(), quarterStart.Month(),quarterStart.Day(),beginTime.Hour(), beginTime.Minute())
        start.AddDate(0,0,firstSectionDay)
        end = time.Date(quarterStart.Year(), quarterStart.Month(),quarterStart.Day(),endTime.Hour(), endTime.Minute())
        end.AddDate(0,0,firstSectionDay)
        
        rule = fmt.Sprintf("FREQ=WEEKLY;UNTIL=%s;BYDAY=%s",quarterEnd.Format("20251205T150000Z"),sectionDays)
        err = createEvent(uid, summary, locaiton, description, rule, start, end, cal)
        if err != nil  {
            return nil, err,i
        }
    }
    return []byte(cal.Serialize()), nil, -1
}

var dayMap = map[rune]string  {
    "M": "MO",
    "T": "TU",
    "W": "WE",
    "R": "TH",
    "F": "FR",
    "S": "SA",
    "U": "SU",
}
func parseDays(s string) string {
    var parts []string
    for _, c := range strings.TrimSpace(s)  {
        if ic, ok := dayMap[c]; ok  {
            parts = append(parts,ic)
        }
    }
    return strings.Join(parts,",")
}
