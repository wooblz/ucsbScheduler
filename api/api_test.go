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
		t.Fatalf("Class slice length mismatch: got %d, want %d", len(got), len(want))
		return // Stop if lengths are different
	}

	for i := range want {
		g_class := got[i]
		w_class := want[i]

		// --- Compare top-level Class fields ---
		if g_class.CourseID != w_class.CourseID {
			t.Errorf("class[%d].CourseID: got %q, want %q", i, g_class.CourseID, w_class.CourseID)
		}
		if g_class.Title != w_class.Title {
			t.Errorf("class[%d].Title: got %q, want %q", i, g_class.Title, w_class.Title)
		}
		if g_class.SubjectArea != w_class.SubjectArea {
			t.Errorf("class[%d].SubjectArea: got %q, want %q", i, g_class.SubjectArea, w_class.SubjectArea)
		}
		if g_class.EnrollCode != w_class.EnrollCode {
			t.Errorf("class[%d].EnrollCode: got %q, want %q", i, g_class.EnrollCode, w_class.EnrollCode)
		}
		if g_class.Room != w_class.Room {
			t.Errorf("class[%d].Room: got %q, want %q", i, g_class.Room, w_class.Room)
		}
		if g_class.Building != w_class.Building {
			t.Errorf("class[%d].Building: got %q, want %q", i, g_class.Building, w_class.Building)
		}
		if g_class.Days != w_class.Days {
			t.Errorf("class[%d].Days: got %q, want %q", i, g_class.Days, w_class.Days)
		}
		if g_class.BeginTime != w_class.BeginTime {
			t.Errorf("class[%d].BeginTime: got %q, want %q", i, g_class.BeginTime, w_class.BeginTime)
		}
		if g_class.EndTime != w_class.EndTime {
			t.Errorf("class[%d].EndTime: got %q, want %q", i, g_class.EndTime, w_class.EndTime)
		}

		// --- Compare ClassSections slice ---
		if len(g_class.ClassSections) != len(w_class.ClassSections) {
			t.Errorf("class[%d] section count mismatch: got %d, want %d", i, len(g_class.ClassSections), len(w_class.ClassSections))
			continue // Skip to next class if section counts differ
		}

		for j := range w_class.ClassSections {
			g_sec := g_class.ClassSections[j]
			w_sec := w_class.ClassSections[j]

			// --- Compare Section fields ---
			if g_sec.Number != w_sec.Number {
				t.Errorf("class[%d].section[%d].Number: got %q, want %q", i, j, g_sec.Number, w_sec.Number)
			}
			if g_sec.EnrollCode != w_sec.EnrollCode {
				t.Errorf("class[%d].section[%d].EnrollCode: got %q, want %q", i, j, g_sec.EnrollCode, w_sec.EnrollCode)
			}

			// --- Compare TimeLocations slice (using DeepEqual is fine here) ---
			if len(g_sec.TimeLocations) != len(w_sec.TimeLocations) {
				t.Errorf("class[%d].section[%d] time location count mismatch: got %d, want %d",
					i, j, len(g_sec.TimeLocations), len(w_sec.TimeLocations))
				continue
			}

			for k := range w_sec.TimeLocations {
				g_loc := g_sec.TimeLocations[k]
				w_loc := w_sec.TimeLocations[k]

				// Since TimeLocation fields should be exact, DeepEqual is still a good choice.
				// For better diffs on a mismatch, consider the go-cmp package.
				if !reflect.DeepEqual(g_loc, w_loc) {
					t.Errorf("class[%d].section[%d].timeLocation[%d] mismatch:\ngot:  %+v\nwant: %+v", i, j, k, g_loc, w_loc)
				}
			}
		}
	}
}

