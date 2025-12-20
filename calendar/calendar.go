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
	"M": time.Monday,
	"T": time.Tuesday,
	"W": time.Wednesday,
	"R": time.Thursday,
	"F": time.Friday,
	"S": time.Saturday,
	"U": time.Sunday,
}

func createEvent(uid, summary, location, description, rule string, start, end time.Time, cal *ics.Calendar) {
	event := cal.AddEvent(uid)
	event.SetSummary(summary)
	event.SetLocation(location)
	event.SetDescription(description)
	event.SetStartAt(start)
	event.SetEndAt(end)
	if rule != "" {
		event.AddRrule(rule)
	}
}

func GenICS(classes []models.Class, finals map[string]models.Final, quarterStart, quarterEnd time.Time) ([]byte, error, int) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)

	for i, c := range classes {
		uid := fmt.Sprintf("%s@ucsbCalendar.com", uuid.NewString())
		courseID := strings.Join(strings.Fields(c.CourseID), " ")
		summary := fmt.Sprintf("%s: %s", courseID, c.Title)
		location := fmt.Sprintf("%s-%s", c.Building, c.Room)
		description := c.SubjectArea

		beginTime, err := time.Parse("15:04", c.BeginTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing begin time for %s: %v", c.CourseID, err), i
		}
		endTime, err := time.Parse("15:04", c.EndTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing end time for %s: %v", c.CourseID, err), i
		}

		firstDayChar := getFirstDayChar(c.Days)
		startDate := calcStartDate(quarterStart, firstDayChar)

		start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), beginTime.Hour(), beginTime.Minute(), 0, 0, time.Local)
		end := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

		classDays := parseDays(c.Days)
		rule := fmt.Sprintf("FREQ=WEEKLY;UNTIL=%s;BYDAY=%s", quarterEnd.Format("20060102T150405Z"), classDays)

		createEvent(uid, summary, location, description, rule, start, end, cal)

		if finals != nil {
			if f, ok := finals[c.EnrollCode]; ok && f.HasFinals {
				finalUid := fmt.Sprintf("FINAL-%s@ucsbCalendar.com", uuid.NewString())
				finalSummary := fmt.Sprintf("Final-%s", summary)

				examDay, err := time.Parse("20060102", f.ExamDate)
				if err != nil {
					return nil, fmt.Errorf("error parsing exam date for %s: %v", c.CourseID, err), i
				}
				examStart, err := time.Parse("15:04", f.BeginTime)
				if err != nil {
					return nil, fmt.Errorf("error parsing exam start for %s: %v", c.CourseID, err), i
				}
				examEnd, err := time.Parse("15:04", f.EndTime)
				if err != nil {
					return nil, fmt.Errorf("error parsing exam end for %s: %v", c.CourseID, err), i
				}

				fStart := time.Date(examDay.Year(), examDay.Month(), examDay.Day(), examStart.Hour(), examStart.Minute(), 0, 0, time.Local)
				fEnd := time.Date(examDay.Year(), examDay.Month(), examDay.Day(), examEnd.Hour(), examEnd.Minute(), 0, 0, time.Local)

				createEvent(finalUid, finalSummary, location, "Final Exam", "", fStart, fEnd, cal)
			}
		}

		if len(c.ClassSections) > 0 && len(c.ClassSections[0].TimeLocations) > 0 {
			t := c.ClassSections[0].TimeLocations[0]

			secUid := fmt.Sprintf("%s-SEC@ucsbCalendar.com", uuid.NewString())
			secSummary := fmt.Sprintf("Section-%s: %s", courseID, c.Title)
			secLocation := fmt.Sprintf("%s-%s", t.Building, t.Room)
			secDescription := "Section"

			secBegin, err := time.Parse("15:04", t.BeginTime)
			if err != nil {
				return nil, fmt.Errorf("error parsing section begin time: %v", err), i
			}
			secEnd, err := time.Parse("15:04", t.EndTime)
			if err != nil {
				return nil, fmt.Errorf("error parsing section end time: %v", err), i
			}

			secFirstDayChar := getFirstDayChar(t.Days)
			secStartDate := calcStartDate(quarterStart, secFirstDayChar)

			secStart := time.Date(secStartDate.Year(), secStartDate.Month(), secStartDate.Day(), secBegin.Hour(), secBegin.Minute(), 0, 0, time.Local)
			secEndT := time.Date(secStartDate.Year(), secStartDate.Month(), secStartDate.Day(), secEnd.Hour(), secEnd.Minute(), 0, 0, time.Local)

			secDays := parseDays(t.Days)
			secRule := fmt.Sprintf("FREQ=WEEKLY;UNTIL=%s;BYDAY=%s", quarterEnd.Format("20060102T150405Z"), secDays)

			createEvent(secUid, secSummary, secLocation, secDescription, secRule, secStart, secEndT, cal)
		}
	}

	return []byte(cal.Serialize()), nil, -1
}

func calcStartDate(quarterStart time.Time, dayChar string) time.Time {
	targetDay, ok := weekDayMap[dayChar]
	if !ok {
		return quarterStart
	}
	daysToAdd := (int(targetDay) - int(quarterStart.Weekday()) + 7) % 7
	return quarterStart.AddDate(0, 0, daysToAdd)
}

func getFirstDayChar(daysStr string) string {
	daysStr = strings.TrimSpace(daysStr)
	if len(daysStr) > 0 {
		return string(daysStr[0])
	}
	return "M"
}

func parseDays(s string) string {
	var parts []string
	for _, r := range strings.ReplaceAll(s, " ", "") {
		c := string(r)
		if ic, ok := iCalDayMap[c]; ok {
			parts = append(parts, ic)
		}
	}
	return strings.Join(parts, ",")
}
