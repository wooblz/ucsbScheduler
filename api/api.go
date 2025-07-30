package api

import (
    "net/http"
    "net/url"
    "strconv"
    "github.com/wooblz/ucsbScheduler/models.go"
    "github.com/joho/godotenv"
)
const loadSize = 500  
func GetAllCourses(quarter int) ([]Class, error) {
    err := godotenv.Load()
    if err != nil   {
        return nil, err
    }
    api_key := os.Getenv("API_KEY")

    baseURL := "https://api.ucsb.edu/academics/curriculums/v3/classes/search"

    parameters := url.Values{}
    parameters.Add("quarter",strconv.Itoa(quarter))
    parameters.Add("pageSize",strconv.Itoa(loadSize))
    parameters.Add("includeClassSections", "true")

    pageCount := 1
    var sol []Class
    for {
        parameter.Set("pageNumber", strconv.Itoa(pageCount))
        url := baseURL + "?" + parameters.Encode()
        client := &http.Client()
        req, err := http.NewRequest("GET", url, nil)
        if err != nil  {
            return nil, err 
        }
        req.Header.Set("Accept", "application/json")
        req.Header.Set("UCSB-API-version", "3.0")
        req.Header.Set("UCSB-API-key", api_key)

        response, err := client.Do(req)
        if err != nil {
            return nil, err
        }
        defer response.Body.Close()
        var result partialApi
        err = json.NewDecoder(response.Body).Decode(&response)
        if err != nil  {
            return nil, err
        }
        sol.Append(sol,result.Classes)
        pageCount++
        if result.Total < loadSize {
            break
        }
        return sol, nil
    }
}
