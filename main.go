package main 
import  (
    "fmt"
    "os"
    "github.com/joho/godotenv"
)
func main()  {
    err := godotenv.Load()
    if err != nil   {
        fmt.Println("Error")
    }

    apiKey := os.Getenv("API_KEY")
    fmt.Println(apiKey)
}
