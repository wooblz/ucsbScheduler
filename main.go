package main

import (
	"database/sql"
	"encoding/json"
//	"fmt"
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

	db, err := database.StartDB("DB_URL")
	if err != nil {
		log.Fatal("Could not connect to DB: ", err)
	}
	defer db.Close()

	err = database.CreateTable(db)
	if err != nil {
		log.Fatal("Could not create tables: ", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM classes").Scan(&count)
	if err != nil {
		log.Println("Error checking database count:", err)
	}

    /*
	fmt.Print("Do you want to repopulate the database? (Y/N): ")
	var input string
	fmt.Scanln(&input)
    
	if count == 0 || input == "Y" || input == "y" {
		refreshDatabase(db)
    }*/

    if count == 0  {
		refreshDatabase(db)
	}

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Println("Starting scheduled 24-hour database refresh...")
				refreshDatabase(db)
				log.Println("Database refresh completed.")
			}
		}
	}()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		handleSearch(w, r, db)
	})

	http.HandleFunc("/api/calendar", func(w http.ResponseWriter, r *http.Request) {
		handleCalendar(w, r, db)
	})

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s...", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func refreshDatabase(db *sql.DB) {
	log.Println("Clearing existing database...")
	err := database.ResetDB(db)
	if err != nil {
		log.Println("Failed to reset database: ", err)
		return
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("API_KEY not set in .env")
		return
	}

	client := &http.Client{Timeout: 60 * time.Second}
	courses, err := api.GetAllCourses(20261, client, "https://api.ucsb.edu/academics/curriculums/v3/classes/search")
	if err != nil {
		log.Println("Failed to fetch courses: ", err)
		return
	}

	err = database.InsertAllClasses(courses, db)
	if err != nil {
		log.Println("Failed to insert courses: ", err)
		return
	}
	log.Println("Database successfully populated")
}

func handleSearch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := r.URL.Query().Get("q")
	log.Printf("Search Request: %s", query)
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}
	results, err := database.QueryTitle(query, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func handleCalendar(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	format := r.URL.Query().Get("format")

	var req CalendarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Calendar Request: Quarter=%s, Codes=%v, Format=%s", req.Quarter, req.EnrollCodes, format)

	quarterInt, err := strconv.Atoi(req.Quarter)
	if err != nil {
		quarterInt = 20261
	}

	var selectedClasses []models.Class
	finalsMap := make(map[string]models.Final)
	selectedCodeMap := make(map[string]bool) 

	client := &http.Client{}
	baseURL := "https://api.ucsb.edu/academics/curriculums/v3/finals"

	for _, code := range req.EnrollCodes {
		selectedCodeMap[code] = true 

		cls, err := database.GetSelectedClass(code, db)
		if err != nil {
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

	start, instrEnd, qEnd := getQuarterDates()

	if format == "preview" {
		instrEnd = start.AddDate(0, 0, 7)
		events, err := calendar.GenerateEvents(selectedClasses, finalsMap, start, instrEnd, qEnd, selectedCodeMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	} else {
		// Pass selectedCodeMap
		icsData, err := calendar.GenICS(selectedClasses, finalsMap, start, instrEnd, qEnd, selectedCodeMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/calendar")
		w.Header().Set("Content-Disposition", "attachment; filename=schedule.ics")
		w.Write(icsData)
	}
}

func getQuarterDates() (time.Time, time.Time, time.Time) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Println("Warning: Could not load America/Los_Angeles location, defaulting to Local")
		loc = time.Local
	}

	start := time.Date(2026, 1, 5, 0, 0, 0, 0, loc)
	instrEnd := time.Date(2026, 3, 13, 23, 59, 59, 0, loc)
	qEnd := time.Date(2026, 3, 20, 23, 59, 59, 0, loc)
	return start, instrEnd, qEnd
}
