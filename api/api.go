package api

import (
    "net/http"
    "net/url"
    "strconv"
    "encoding/json"
    "os"
    "math"
    "errors"
    "fmt"
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
        if result.Message != ""  {
            return nil, errors.New(result.Message)
        }
        if result.Total == 0  {
            break
        }
        for j := range result.Classes  {
            c := result.Classes[j]
            classList, err := getMainTime(c)
            if err != nil  {
                continue 
            }
            sol = append(sol, classList...)
        }
        pageCount++
        if math.Ceil(float64(result.Total)/float64(loadSize)) < float64(pageCount) {
            break
        }
    }
    return sol, nil
}
func getMainTime(c models.Class) ([]models.Class, error)  {
    var classes []models.Class
    dic := make(map[string][]models.Section)
    for _, s := range c.ClassSections {
        firstTwo := s.Number[:2]
        dic[firstTwo] = append(dic[firstTwo], s)
    }
    for key, value := range dic  {

        cur := models.Class  {
            CourseID: c.CourseID,
            Title: c.Title,
            SubjectArea: c.SubjectArea,
        }
        a,_ := getMainSection(key, value)

        cur.ClassSections = value 
        s := value[a].TimeLocations[0]
        cur.Room = s.Room
        cur.Building = s.Building
        cur.Days = s.Days
        cur.BeginTime = s.BeginTime
        cur.EndTime = s.EndTime
        cur.EnrollCode = value[a].EnrollCode
        cur.ClassSections = append(cur.ClassSections[:a], cur.ClassSections[a+1:]...)
        classes = append(classes, cur)
    }
    return classes, nil

}
func getMainSection(code string, sections []models.Section) (int, error)  {
    for i, s := range sections  {
        if s.Number == code + "00"  {
            return i, nil
        }
    }
    return -1, fmt.Errorf("Failed to find main class")
}

//https://api.ucsb.edu/academics/curriculums/v3/finals
func GetFinal(quarter int, enrollCode string, baseURL string, client *http.Client) (models.Final,error) {
    if len(strconv.Itoa(quarter)) != 5  {
        return models.Final{}, errors.New("Invalid quarter, YYYYQ format")
    }
    err := godotenv.Load("../.env")
    if err != nil   {
        return models.Final{}, err
    }
    api_key := os.Getenv("API_KEY")
    parameters := url.Values{}
    parameters.Add("quarter",strconv.Itoa(quarter))
    parameters.Add("enrollCode", enrollCode)

    url := baseURL + "?" + parameters.Encode()
    req,err := http.NewRequest("GET", url, nil)
    if err !=  nil  {
        return models.Final{}, err
    }

    req.Header.Set("Accept", "application/json")
    req.Header.Set("UCSB-API-version", "3.0")
    req.Header.Set("UCSB-API-key", api_key)

    response, err := client.Do(req)
    if err != nil {
        return models.Final{}, err
    }
    defer response.Body.Close()
    var result models.Final
        
    err = json.NewDecoder(response.Body).Decode(&result)
    if err != nil  {
        return models.Final{}, err
    }
    if result.Message != ""  {
        return result, errors.New(result.Message)
    }
    if result.ExamDay == ""  {
        return result, errors.New(fmt.Sprintf("Could not find final for enrollCode: %s", enrollCode))
    }
    return result, nil
}
