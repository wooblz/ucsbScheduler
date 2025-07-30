package api

import  (
    "testing"
    "net/http/httptest"
    "os"
    "net/http"
    "io"
    //"github.com/wooblz/ucsbScheduler/models"
    
)

func TestAPI(t *testing.T)  {
    t.Run("test single input", func (t *testing.T)  {
        data,err := os.ReadFile("api_test_data/oneclassandsection")
        if err != nil   {
            t.Fatalf("Failed to open file: %v", err)
        }
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            io.WriteString(w, string(data))
        }))
        defer server.Close()
        client := server.Client()
        got, err := GetAllCourses(20251, client, server.URL)
        if err != nil  {
            t.Fatalf("Failed to get course: %v", err)
        }
        class := got[0]
        assertEqual(t, "CourseID", class.CourseID, "CMPSC     5A ")
        assertEqual(t, "Title", class.Title, "INTRO DATA SCI 1")
        assertEqual(t, "SubjectArea", class.SubjectArea, "CMPSC   ")

        if len(class.ClassSections) != 1 {
            t.Fatalf("unexpected ClassSections length: got %d, want 1", len(class.ClassSections))
        }

        section := class.ClassSections[0]

        if len(section.TimeLocations) != 1 {
            t.Fatalf("unexpected TimeLocations length: got %d, want 1", len(section.TimeLocations))
        }

        timeLoc := section.TimeLocations[0]

        assertEqual(t, "Room", timeLoc.Room, "1610")
        assertEqual(t, "Building", timeLoc.Building, "BRDA")
        assertEqual(t, "Days", timeLoc.Days, "M W    ")
        assertEqual(t, "BeginTime", timeLoc.BeginTime, "15:30")
        assertEqual(t, "EndTime", timeLoc.EndTime, "16:45")
    })
}
func assertEqual(t *testing.T, name string, got, want string) {
    if got != want {
        t.Fatalf("unexpected %s: got %q, want %q", name, got, want)
    }
}
