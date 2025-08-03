package api

import (
    "net/http"
    "net/url"
    "strconv"
    "encoding/json"
    "os"
    "math"
    "errors"
    "github.com/wooblz/ucsbScheduler/models"
    "github.com/joho/godotenv"
)
var loadSize = 500
//"https://api.ucsb.edu/academics/curriculums/v3/classes/search"
func GetAllCourses(quarter int, client *http.Client, baseURL string) ([]models.Class, error) {
    if len(strconv.Itoa(quarter)) != 5  {
        return nil, errors.New("Invalid quarter, YYYYQ format")
    }
    err := godotenv.Load("../.env")
    if err != nil   {
        return nil, err
    }
    api_key := os.Getenv("API_KEY")


    parameters := url.Values{}
    parameters.Add("quarter",strconv.Itoa(quarter))
    parameters.Add("pageSize",strconv.Itoa(loadSize))
    parameters.Add("includeClassSections", "true")

    pageCount := 1
    var sol []models.Class
    for {
        parameters.Set("pageNumber", strconv.Itoa(pageCount))
        url := baseURL + "?" + parameters.Encode()
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
        var result models.PartialApi
        err = json.NewDecoder(response.Body).Decode(&result)
        if err != nil  {
            return nil, err
        }
        if result.Total == 0  {
            break
        }
        if result.Message != ""  {
            return nil, errors.New(result.Message)
        }
        sol = append(sol,result.Classes...)
        pageCount++
        if math.Ceil(float64(result.Total)/float64(loadSize)) < float64(pageCount) {
            break
        }
    }
    return sol, nil
}
