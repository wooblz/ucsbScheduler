package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wooblz/ucsbScheduler/api"
	"github.com/wooblz/ucsbScheduler/calendar"
	"github.com/wooblz/ucsbScheduler/database"
	"github.com/wooblz/ucsbScheduler/models"
)

type CalendarRequest struct {
	Quarter     string   `json:"quarter"`
	EnrollCodes []string `json:"enrollCodes"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// 1. Connect to Database
	db, err := database.StartDB("DB_URL")
	if err != nil {
		log.Fatal("Could not connect to DB: ", err)
	}
	defer db.Close()

	// 2. Setup Tables
	err = database.CreateTable(db)
	if err != nil {
		log.Fatal("Could not create tables: ", err)
	}

	// ==========================================
	// DEBUG: Test Database Content
	// ==========================================
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM classes").Scan(&count)
	if err != nil {
		log.Println("Error checking database count:", err)
	}
	log.Printf("DEBUG: Database currently contains %d classes.", count)

	// 3. Populate DB if empty
	if count == 0 {
		log.Println("⚠️ Database is empty. Attempting to fetch courses from UCSB API...")
		
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			log.Fatal("API_KEY not set in .env. Cannot fetch courses.")
		}

		client := &http.Client{Timeout: 60 * time.Second}
		// Fetching Winter 2025 (20251)
		courses, err := api.GetAllCourses(20251, client, "https://api.ucsb.edu/academics/curriculums/v3/classes/search")
		if err != nil {
			log.Fatal("Failed to fetch courses from API: ", err)
		}
		
		log.Printf("Fetched %d courses. Inserting into database...", len(courses))
		err = database.InsertAllClasses(courses, db)
		if err != nil {
			log.Fatal("Failed to insert courses: ", err)
		}
		log.Println("✅ Database successfully populated!")
	}

	// 4. Test a specific query in Console
	testQuery := "CMPSC"
	log.Printf("DEBUG: Running test search for '%s'...", testQuery)
	results, err := database.QueryTitle(testQuery, db)
	if err != nil {
		log.Printf("DEBUG: Search error: %v", err)
	} else {
		log.Printf("DEBUG: Found %d results for '%s'. Top result: %v", len(results), testQuery, 
			func() string { if len(results) > 0 { return results[0].Title } else { return "None" } }())
	}
	// ==========================================

	// 5. Start Server
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		handleSearch(w, r, db)
	})

	http.HandleFunc("/api/calendar", func(w http.ResponseWriter, r *http.Request) {
		handleCalendar(w, r, db)
	})

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSearch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	results, err := database.QueryTitle(query, db)
	if err != nil {
		log.Printf("Search Error: %v", err) // Log error to console
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Debug: Log what is being sent to frontend
	log.Printf("Query: %s | Results found: %d", query, len(results))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func handleCalendar(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req CalendarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	baseURL := "https://api.ucsb.edu/academics/curriculums/v3/finals"

	quarterInt, err := strconv.Atoi(req.Quarter)
	if err != nil {
		http.Error(w, "Invalid quarter format", http.StatusBadRequest)
		return
	}

	var selectedClasses []models.Class
	finalsMap := make(map[string]models.Final)

	for _, code := range req.EnrollCodes {
		cls, err := database.GetSelectedClass(code, db)
		if err != nil {
			log.Printf("Warning: Could not find class for enroll code %s: %v", code, err)
			continue
		}
		selectedClasses = append(selectedClasses, cls)

		if cls.EnrollCode != "" {
			final, err := api.GetFinal(quarterInt, cls.EnrollCode, baseURL, client)
			if err == nil {
				finalsMap[cls.EnrollCode] = final
			}
		}
	}

	start, end := getQuarterDates(req.Quarter)

	icsData, err, _ := calendar.GenICS(selectedClasses, finalsMap, start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", "attachment; filename=schedule.ics")
	w.Write(icsData)
}

func getQuarterDates(quarter string) (time.Time, time.Time) {
	// Updated to requested start date: Jan 5, 2025
	if quarter == "20251" {
		return time.Date(2025, 1, 5, 0, 0, 0, 0, time.Local), time.Date(2025, 3, 14, 23, 59, 59, 0, time.Local)
	}
	return time.Now(), time.Now()
}
