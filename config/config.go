package config

import "time"

func GetQuarterDates(quarter string) (time.Time, time.Time) {
    if quarter == "20251" { 
        start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.Local)
        end := time.Date(2025, 3, 14, 0, 0, 0, 0, time.Local)
        return start, end
    }
    return time.Now(), time.Now()
}
