package api

import  (
    "testing"
    "net/http/httptest"
    "os"
    "net/http"
    "io"
    "reflect"
    "github.com/wooblz/ucsbScheduler/models"
    
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

        assertClassesEqual(t, got, Solution1)
    })
    t.Run("test multiple classes", func(t *testing.T)  {
        data,err := os.ReadFile("api_test_data/multiclass")
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
        assertClassesEqual(t, got, Solution2)
    })

}
func assertClassesEqual(t *testing.T, got, want []models.Class) {
    if len(got) != len(want) {
        t.Fatalf("Class count mismatch: got %d, want %d", len(got), len(want))
    }

    for i := range want {
        if got[i].CourseID != want[i].CourseID {
            t.Errorf("class[%d].CourseID: got %q, want %q", i, got[i].CourseID, want[i].CourseID)
        }
        if got[i].Title != want[i].Title {
            t.Errorf("class[%d].Title: got %q, want %q", i, got[i].Title, want[i].Title)
        }
        if got[i].SubjectArea != want[i].SubjectArea {
            t.Errorf("class[%d].SubjectArea: got %q, want %q", i, got[i].SubjectArea, want[i].SubjectArea)
        }
        if len(got[i].ClassSections) != len(want[i].ClassSections) {
            t.Errorf("class[%d] section count mismatch: got %d, want %d", i, len(got[i].ClassSections), len(want[i].ClassSections))
            continue
        }

        for j := range want[i].ClassSections {
            if len(got[i].ClassSections[j].TimeLocations) != len(want[i].ClassSections[j].TimeLocations) {
                t.Errorf("class[%d] section[%d] time location count mismatch: got %d, want %d",
                    i, j, len(got[i].ClassSections[j].TimeLocations), len(want[i].ClassSections[j].TimeLocations))
                continue
            }

            for k := range want[i].ClassSections[j].TimeLocations {
                g := got[i].ClassSections[j].TimeLocations[k]
                w := want[i].ClassSections[j].TimeLocations[k]
                if !reflect.DeepEqual(g, w) {
                    t.Errorf("class[%d] section[%d] timeLocation[%d] mismatch:\ngot:  %+v\nwant: %+v", i, j, k, g, w)
                }
            }
        }
    }
}

