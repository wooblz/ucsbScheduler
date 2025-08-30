package api


import  (
    "testing"
    "net/http/httptest"
    "os"
    "net/http"
    "io"
    "reflect"
    "strconv"
    "github.com/wooblz/ucsbScheduler/models"
    
)
func TestFinalAPI(t *testing.T)  {
    t.Run("final", func(t *testing.T)  {
        data,err  := os.ReadFile("api_test_data/final.txt")
        if err != nil  {
            t.Fatalf("Failed to open file: %v", err)  
        }
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
            w.Header().Set("Content-Type", "application/json")
            io.WriteString(w, string(data))
        }))
        defer server.Close()
        client := server.Client()
        got, err := GetFinal(20251,"test",server.URL, client)
        if err != nil  {
            t.Fatalf("Failed to get final: %v", err)
        }
        assertFinalEqual(t, got, FinalSolution1)
    })
    t.Run("can't find final", func(t *testing.T)  {
        data,err  := os.ReadFile("api_test_data/finalnull")
        if err != nil  {
            t.Fatalf("Failed to open file: %v", err)  
        }
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
            w.Header().Set("Content-Type", "application/json")
            io.WriteString(w, string(data))
        }))
        defer server.Close()
        client := server.Client()
        _, err = GetFinal(20251,"test",server.URL, client)
        if err == nil  {
            t.Fatalf("Failed to cause error: %v", err)
        }
    })
}
func TestClassAPI(t *testing.T)  {
    /*t.Run("test single input", func (t *testing.T)  {
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
    })*/
    t.Run("multi page test", func(t *testing.T)  {
        old := loadSize
        loadSize = 3
        defer func() {loadSize = old}()

        
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            pageNumber := r.URL.Query().Get("pageNumber")
            page,err := strconv.Atoi(pageNumber)
            if err != nil  {
                t.Fatalf("Invalid page number %v", err)
            }
            switch page  {
                case 1:
                    data,err := os.ReadFile("api_test_data/page1.txt")
                    if err != nil   {
                        t.Fatalf("Failed to open file: %v", err)
                    }
                    io.WriteString(w, string(data))
                case 2:
                    data,err := os.ReadFile("api_test_data/page2.txt")
                    if err != nil   {
                        t.Fatalf("Failed to open file: %v", err)
                    }
                    io.WriteString(w, string(data))
            }
        }))
        defer server.Close()
        client := server.Client()
        got, err := GetAllCourses(20251, client, server.URL)
        if err != nil  {
            t.Fatalf("Failed to get course: %v", err)
        }
        assertClassesEqual(t, got, Solution3)
    })
    t.Run("empty query", func(t *testing.T)  {
        data, err := os.ReadFile("api_test_data/emptyquery")
        if err != nil  {
            t.Fatalf("Failed to open file: %v", err)
        }
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
            w.Header().Set("Content-Type", "application/json")
            io.WriteString(w, string(data))
        }))
        defer server.Close()
        client := server.Client()
        got, err := GetAllCourses(20251, client, server.URL)
        if err  != nil  {
            t.Fatalf("Failed to get courses: %v", err)
        }
        if len(got) != 0  {
            t.Errorf("Empty Query Failed")
        }
    })
    t.Run("error query", func(t *testing.T)  {
        data, err := os.ReadFile("api_test_data/errorquery")
        if err != nil  {
            t.Fatalf("Failed to open file: %v", err)
        }
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)  {
            w.Header().Set("Content-Type", "application/json")
            io.WriteString(w, string(data))
        }))
        defer server.Close()
        client := server.Client()
        _, err = GetAllCourses(20251, client, server.URL)
        if err == nil  {
            t.Fatalf("Failed to cause error")
        }
    })

}
func assertFinalEqual(t *testing.T, got, want models.Final)  {
    if got.HasFinals != want.HasFinals {
        t.Errorf("HasFinals mismatch: got %t, want %t", got.HasFinals, want.HasFinals)
    }
    if got.Comments != want.Comments {
        t.Errorf("Comments mismatch: got %q, want %q", got.Comments, want.Comments)
    }
    if got.ExamDay != want.ExamDay {
        t.Errorf("ExamDay mismatch: got %q, want %q", got.ExamDay, want.ExamDay)
    }
    if got.ExamDate != want.ExamDate {
        t.Errorf("ExamDate mismatch: got %q, want %q", got.ExamDate, want.ExamDate)
    }
    if got.BeginTime != want.BeginTime {
        t.Errorf("BeginTime mismatch: got %q, want %q", got.BeginTime, want.BeginTime)
    }
    if got.EndTime != want.EndTime {
        t.Errorf("EndTime mismatch: got %q, want %q", got.EndTime, want.EndTime)
    }
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

