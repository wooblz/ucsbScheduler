package calendar
import  (
    "time"
)
func generateCalendarBytes() ([]byte, error) {
	// 1. Create a new calendar
	cal := ics.NewCalendar()

	// --- Event 1: CMPSC 16 Lecture ---
	event1 := cal.AddEvent("test-uid-1@ucsbCalendar.com")
	event1.SetSummary("CMPSC 16: PROBLEM SOLVING I")
	event1.SetLocation("ILP-1101")
	event1.SetDescription("CMPSC")

	// Set start and end times for the first occurrence
	start1, err := time.Parse(time.RFC3339, "2025-09-29T14:00:00Z")
	if err != nil {
		return nil, fmt.Errorf("error parsing start time for event 1: %w", err)
	}
	end1, err := time.Parse(time.RFC3339, "2025-09-29T15:15:00Z")
	if err != nil {
		return nil, fmt.Errorf("error parsing end time for event 1: %w", err)
	}
	event1.SetStartAt(start1)
	event1.SetEndAt(end1)

	// Add the recurrence rule
	event1.AddRrule("FREQ=WEEKLY;UNTIL=20251205T235959Z;BYDAY=MO,WE")

	// --- Event 2: CMPSC 16 Section ---
	event2 := cal.AddEvent("test-uid-2@ucsbCalendar.com")
	event2.SetSummary("Section-CMPSC 16: PROBLEM SOLVING I")
	event2.SetLocation("PHELP-3525")
	event2.SetDescription("Section")

	start2, err := time.Parse(time.RFC3339, "2025-09-25T09:00:00Z")
	if err != nil {
		return nil, fmt.Errorf("error parsing start time for event 2: %w", err)
	}
	end2, err := time.Parse(time.RFC3339, "2025-09-25T09:50:00Z")
	if err != nil {
		return nil, fmt.Errorf("error parsing end time for event 2: %w", err)
	}
	event2.SetStartAt(start2)
	event2.SetEndAt(end2)

	event2.AddRrule("FREQ=WEEKLY;UNTIL=20251205T235959Z;BYDAY=TH")

	// --- Event 3: MATH 190PS Lecture ---
	event3 := cal.AddEvent("test-uid-3@ucsbCalendar.com")
	event3.SetSummary("MATH 190PS: PROBLEM SOLVING")
	event3.SetLocation("PHELP-1508")
	event3.SetDescription("MATH")

	start3, err := time.Parse(time.RFC3339, "2025-09-29T09:00:00Z")
	if err != nil {
		return nil, fmt.Errorf("error parsing start time for event 3: %w", err)
	}
	end3, err := time.Parse(time.RFC3339, "2025-09-29T09:50:00Z")
	if err != nil {
		return nil, fmt.Errorf("error parsing end time for event 3: %w", err)
	}
	event3.SetStartAt(start3)
	event3.SetEndAt(end3)

	event3.AddRrule("FREQ=WEEKLY;UNTIL=20251205T235959Z;BYDAY=MO,WE,FR")

	return []byte(cal.Serialize()), nil
}
var Solution3 = []models.Class{
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
	cal.SetVersion("2.0")
    {
        CourseID:    "MATH    190PS",
        Title:       "PROBLEM SOLVING",
        SubjectArea: "MATH    ",
        Room:        "1508",
        Building:    "PHELP",
        Days:        "M W F  ",
        BeginTime:   "09:00",
        EndTime:     "09:50",
        ClassSections: []models.Section{},
    },
}
