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

	db, err := database.StartDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	if quarter == "20251" {
		return time.Date(2025, 1, 6, 0, 0, 0, 0, time.Local), time.Date(2025, 3, 14, 0, 0, 0, 0, time.Local)
	}
	return time.Now(), time.Now()
}
