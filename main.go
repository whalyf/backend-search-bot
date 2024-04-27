package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"
		"fmt"
		"os"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
		"github.com/joho/godotenv"
		// g "github.com/serpapi/google-search-results-golang"
)

func main() {
    router := mux.NewRouter()

		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

    // Register the handler for the "/process" endpoint with Gorilla Mux
    router.HandleFunc("/process", handleProcessRequest).Methods("POST")

    corsHandler := handlers.CORS(
        handlers.AllowedHeaders([]string{"Content-Type"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        handlers.AllowedOrigins([]string{"*"}), // Allow requests from any origin
    )

    // Wrap your router with the CORS middleware
    http.Handle("/", corsHandler(router))

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleProcessRequest(w http.ResponseWriter, r *http.Request) {
    var searchParams map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&searchParams)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    processedResults := processNestJSData(searchParams)
  	// Check if "keywords" exists in processedResults
		keywords, ok := processedResults["keywords"].(string)
		if !ok {
				http.Error(w, "Keywords not found or not a string", http.StatusInternalServerError)
				return
		}
		email, ok := processedResults["email"].(string)
		if !ok {
				http.Error(w, "Email not found or not a string", http.StatusInternalServerError)
				return
		}

		// Call searchOnGoogle with the keywords
		results:= searchOnGoogle(keywords, email)
		fmt.Println(results)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

func searchOnGoogle(keywords string, email string) map[string]string{
	parameter := map[string]string{
    "api_key": os.Getenv("SERPAPI_KEY"),
    "engine": "google",
    "q": keywords,
    "location": "Brazil",
    "google_domain": "google.com.br",
    "gl": "br",
    "hl": "pt",
  }
  // search := g.NewGoogleSearch(parameter, os.Getenv("SERPAPI_KEY"))
  // results, err := search.GetJSON()

	// fmt.Println(err)
	return parameter
}

func processNestJSData(searchParams map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "location":  searchParams["location"],
        "frequency": searchParams["frequency"],
        "email": searchParams["email"],
        "keywords":  toUpperCase(searchParams["keywords"].(string)),
        "dateTime":  time.Now().Format(time.RFC3339),
        "searchId":   searchParams["searchId"],
    }
}

func toUpperCase(s string) string {
    return strings.ToUpper(s)
}
