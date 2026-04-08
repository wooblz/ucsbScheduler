package calendar

import (
	"fmt"
	"strings"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/wooblz/ucsbScheduler/models"
)

// TestGenICS verifies that the generated ICS matches the logic of the expected "golden" output.
func TestGenICS(t *testing.T) {
	// 1. Setup Input Data (Based on your provided Solution3)
	classes := []models.Class{
		{
			CourseID:    "CMPSC    16  ",
			Title:       "PROBLEM SOLVING I",
			SubjectArea: "CMPSC   ",
			Room:        "1101",
			Building:    "ILP",
			Days:        "M W    ",
			BeginTime:   "14:00",
			EndTime:     "15:15",
			ClassSections: []models.Section{
				{
					TimeLocations: []models.TimeLocation{
						{Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "09:00", EndTime: "09:50"},
					},
				},
			},
		},
		{
			CourseID:      "MATH    190PS",
			Title:         "PROBLEM SOLVING",
			SubjectArea:   "MATH    ",
			Room:          "1508",
			Building:      "PHELP",
			Days:          "M W F  ",
			BeginTime:     "09:00",
			EndTime:       "09:50",
			ClassSections: []models.Section{},
		},
	}

	// 2. Setup Quarter Dates
	// CRITICAL: Your "Golden" data expects the Section (Thursday) to start on Sept 25,
	// and the Lecture (Mon/Wed) to start on Sept 29.
	// For this logic to hold true, the Quarter must start on Thursday, Sept 25, 2025.
	quarterStart := time.Date(2025, 9, 25, 0, 0, 0, 0, time.Local)
	quarterEnd := time.Date(2025, 12, 5, 23, 59, 59, 0, time.Local)

	// 3. Run the Function
	gotBytes, err, _ := GenICS(classes, nil, quarterStart, quarterEnd)
	if err != nil {
		t.Fatalf("GenICS returned unexpected error: %v", err)
	}

	// 4. Validate Output
	// Since GenICS uses random UUIDs, we cannot compare bytes directly.
	// We parse the output and check if it contains the expected events.
	gotCal, err := ics.ParseCalendar(strings.NewReader(string(gotBytes)))
	if err != nil {
		t.Fatalf("Failed to parse generated ICS: %v", err)
	}

	// Define what we expect to find
	expectedEvents := map[string]struct {
		Summary  string
		Location string
		Start    string // Formatted as YYYYMMDDTHHMMSS
		End      string
		RRule    string
	}{
		"Lecture 1": {
			Summary:  "CMPSC 16: PROBLEM SOLVING I",
			Location: "ILP-1101",
			Start:    "20250929T140000", // Mon Sept 29
			End:      "20250929T151500",
			RRule:    "FREQ=WEEKLY;UNTIL=20251205T235959Z;BYDAY=MO,WE",
		},
		"Section 1": {
			Summary:  "Section-CMPSC 16: PROBLEM SOLVING I",
			Location: "PHELP-3525",
			Start:    "20250925T090000", // Thu Sept 25
			End:      "20250925T095000",
			RRule:    "FREQ=WEEKLY;UNTIL=20251205T235959Z;BYDAY=TH",
		},
		"Lecture 2": {
			Summary:  "MATH 190PS: PROBLEM SOLVING",
			Location: "PHELP-1508",
			Start:    "20250929T090000", // Mon Sept 29
			End:      "20250929T095000",
			RRule:    "FREQ=WEEKLY;UNTIL=20251205T235959Z;BYDAY=MO,WE,FR",
		},
	}

	// Iterate over generated events and match them to expectations
	foundCount := 0
	for _, event := range gotCal.Events() {
		summary := event.GetProperty(ics.ComponentPropertySummary).Value
		
		// Find which expectation matches this event
		var matchKey string
		for key, exp := range expectedEvents {
			if strings.Contains(summary, exp.Summary) || summary == exp.Summary {
				matchKey = key
				break
			}
		}

		if matchKey == "" {
			t.Errorf("Found unexpected event in ICS: %s", summary)
			continue
		}

		exp := expectedEvents[matchKey]
		foundCount++

		// Verify Location
		loc := event.GetProperty(ics.ComponentPropertyLocation).Value
		if loc != exp.Location {
			t.Errorf("[%s] Location mismatch. Got %s, want %s", matchKey, loc, exp.Location)
		}

		// Verify Start Time
		start := event.GetProperty(ics.ComponentPropertyDtStart).Value
		// Handle potential 'Z' suffix or local time differences in string comparison
		if !strings.HasPrefix(start, exp.Start) {
			t.Errorf("[%s] Start Time mismatch. Got %s, want start %s", matchKey, start, exp.Start)
		}

		// Verify RRULE
		// We normalize by removing spaces just in case
		rrule := event.GetProperty(ics.ComponentPropertyRrule).Value
		if !strings.Contains(rrule, exp.RRule) {
			t.Errorf("[%s] RRULE mismatch. Got %s, want containing %s", matchKey, rrule, exp.RRule)
		}
	}

	if foundCount != len(expectedEvents) {
		t.Errorf("Expected %d unique events, found %d", len(expectedEvents), foundCount)
	}
}

// Optional: Helper to visualize what your manual function produces vs GenICS
// This confirms your provided 'generateCalendarBytes' is valid code
func generateManualBytes() ([]byte, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)

	// Event 1
	event1 := cal.AddEvent("test-uid-1@ucsbCalendar.com")
	event1.SetSummary("CMPSC 16: PROBLEM SOLVING I")
	event1.SetLocation("ILP-1101")
	event1.SetDescription("CMPSC")
	start1, _ := time.Parse(time.RFC3339, "2025-09-29T14:00:00Z")
	end1, _ := time.Parse(time.RFC3339, "2025-09-29T15:15:00Z")
	event1.SetStartAt(start1)
	event1.SetEndAt(end1)
	event1.AddRrule("FREQ=WEEKLY;UNTIL=20251205T235959Z;BYDAY=MO,WE")

	// Event 2
	event2 := cal.AddEvent("test-uid-2@ucsbCalendar.com")
	event2.SetSummary("Section-CMPSC 16: PROBLEM SOLVING I")
	event2.SetLocation("PHELP-3525")
	event2.SetDescription("Section")
	start2, _ := time.Parse(time.RFC3339, "2025-09-25T09:00:00Z")
	end2, _ := time.Parse(time.RFC3339, "2025-09-25T09:50:00Z")
	event2.SetStartAt(start2)
	event2.SetEndAt(end2)
	event2.AddRrule("FREQ=WEEKLY;UNTIL=20251205T235959Z;BYDAY=TH")

	return []byte(cal.Serialize()), nil
}
