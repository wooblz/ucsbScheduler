package calendar

import (
	"fmt"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
	"github.com/wooblz/ucsbScheduler/models"
)

var iCalDayMap = map[string]string{
	"M": "MO", "T": "TU", "W": "WE", "R": "TH", "F": "FR", "S": "SA", "U": "SU",
}

var weekDayMap = map[string]time.Weekday{
	"M": time.Monday, "T": time.Tuesday, "W": time.Wednesday, "R": time.Thursday, "F": time.Friday, "S": time.Saturday, "U": time.Sunday,
}

type Event struct {
	Title       string `json:"title"`
	Start       string `json:"start"`
	End         string `json:"end,omitempty"`
	Duration    string `json:"duration,omitempty"`
	RRule       string `json:"rrule,omitempty"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Color       string `json:"color,omitempty"`
}

func GenerateEvents(classes []models.Class, finals map[string]models.Final, quarterStart, instructionEnd, quarterEnd time.Time) ([]Event, error) {
	var events []Event

	for _, c := range classes {
		courseID := strings.Join(strings.Fields(c.CourseID), " ")
		summary := fmt.Sprintf("%s: %s", courseID, c.Title)
		location := fmt.Sprintf("%s-%s", c.Building, c.Room)

		if c.Days != "" && c.BeginTime != "" && c.EndTime != "" {
			ev, err := createRecurringEvent(
				summary, location, c.SubjectArea, "#3788d8",
				c.Days, c.BeginTime, c.EndTime,
				quarterStart, instructionEnd,
			)
			if err != nil { return nil, err }
			events = append(events, ev)
		}

		if len(c.ClassSections) > 0 && len(c.ClassSections[0].TimeLocations) > 0 {
			t := c.ClassSections[0].TimeLocations[0]
			secSummary := fmt.Sprintf("Section-%s", courseID)
			secLocation := fmt.Sprintf("%s-%s", t.Building, t.Room)

			if t.Days != "" && t.BeginTime != "" && t.EndTime != "" {
				ev, err := createRecurringEvent(
					secSummary, secLocation, "Section", "#28a745",
					t.Days, t.BeginTime, t.EndTime,
					quarterStart, instructionEnd,
				)
				if err != nil { return nil, err }
				events = append(events, ev)
			}
		}

		if finals != nil {
			if f, ok := finals[c.EnrollCode]; ok && f.HasFinals {
				finalSummary := fmt.Sprintf("Final-%s", courseID)
				
				examDay, err := time.Parse("20060102", f.ExamDate)
				if err != nil { continue }
				examStart, err := time.Parse("15:04", f.BeginTime)
				if err != nil { continue }
				examEnd, err := time.Parse("15:04", f.EndTime)
				if err != nil { continue }

				fStart := time.Date(examDay.Year(), examDay.Month(), examDay.Day(), examStart.Hour(), examStart.Minute(), 0, 0, time.Local)
				fEnd := time.Date(examDay.Year(), examDay.Month(), examDay.Day(), examEnd.Hour(), examEnd.Minute(), 0, 0, time.Local)

				events = append(events, Event{
					Title:       finalSummary,
					Start:       fStart.Format("2006-01-02T15:04:05"),
					End:         fEnd.Format("2006-01-02T15:04:05"),
					Location:    location,
					Description: "Final Exam",
					Color:       "#dc3545",
				})
			}
		}
	}
	return events, nil
}

func createRecurringEvent(title, location, desc, color, days, beginTimeStr, endTimeStr string, qStart, qEnd time.Time) (Event, error) {
	beginTime, err := time.Parse("15:04", beginTimeStr)
	if err != nil { return Event{}, err }
	endTime, err := time.Parse("15:04", endTimeStr)
	if err != nil { return Event{}, err }

	firstDayChar := getFirstDayChar(days)
	startDate := calcStartDate(qStart, firstDayChar)

	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), beginTime.Hour(), beginTime.Minute(), 0, 0, time.Local)
	
	dur := endTime.Sub(beginTime)
	hours := int(dur.Hours())
	minutes := int(dur.Minutes()) % 60
	durationStr := fmt.Sprintf("%02d:%02d", hours, minutes)

	classDays := parseDays(days)
	untilStr := qEnd.UTC().Format("20060102T150405Z")
	rrule := fmt.Sprintf("FREQ=WEEKLY;UNTIL=%s;BYDAY=%s", untilStr, classDays)

	return Event{
		Title:       title,
		Start:       start.Format("2006-01-02T15:04:05"), 
		Duration:    durationStr,
		RRule:       rrule,
		Location:    location,
		Description: desc,
		Color:       color,
	}, nil
}

func GenICS(classes []models.Class, finals map[string]models.Final, quarterStart, instructionEnd, quarterEnd time.Time) ([]byte, error) {
	events, err := GenerateEvents(classes, finals, quarterStart, instructionEnd, quarterEnd)
	if err != nil { return nil, err }

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)

	for _, e := range events {
		uid := fmt.Sprintf("%s@ucsbCalendar.com", uuid.NewString())
		event := cal.AddEvent(uid)
		event.SetSummary(e.Title)
		event.SetLocation(e.Location)
		event.SetDescription(e.Description)
		
		s, _ := time.ParseInLocation("2006-01-02T15:04:05", e.Start, time.Local)
		event.SetStartAt(s)
		
		if e.RRule != "" {
			event.AddRrule(e.RRule)
			d, _ := time.ParseDuration(strings.Replace(e.Duration, ":", "h", 1) + "m")
			event.SetEndAt(s.Add(d))
		} else {
			end, _ := time.ParseInLocation("2006-01-02T15:04:05", e.End, time.Local)
			event.SetEndAt(end)
		}
	}
	return []byte(cal.Serialize()), nil
}

func calcStartDate(quarterStart time.Time, dayChar string) time.Time {
	targetDay, ok := weekDayMap[dayChar]
	if !ok { return quarterStart }
	daysToAdd := (int(targetDay) - int(quarterStart.Weekday()) + 7) % 7
	return quarterStart.AddDate(0, 0, daysToAdd)
}

func getFirstDayChar(daysStr string) string {
	daysStr = strings.TrimSpace(daysStr)
	if len(daysStr) > 0 { return string(daysStr[0]) }
	return "M"
}

func parseDays(s string) string {
	var parts []string
	for _, r := range strings.ReplaceAll(s, " ", "") {
		c := string(r)
		if ic, ok := iCalDayMap[c]; ok { parts = append(parts, ic) }
	}
	return strings.Join(parts, ",")
}
