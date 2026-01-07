package config

import "time"

    
    if quarter == "20261" {
        start := time.Date(2026, 1, 5, 0, 0, 0, 0, time.Local)
        end := time.Date(2026, 3, 20, 23, 59, 59, 0, time.Local)
        return start, end, nil
    }
    
    return time.Time{}, time.Time{}, fmt.Errorf("quarter %s not defined", quarter)
}
