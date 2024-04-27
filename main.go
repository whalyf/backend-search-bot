package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()

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

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(processedResults)
}

func processNestJSData(searchParams map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "location":  searchParams["location"],
        "frequency": searchParams["frequency"],
        "keywords":  toUpperCase(searchParams["keywords"].(string)),
        "dateTime":  time.Now().Format(time.RFC3339),
        "searchId":   searchParams["searchId"],
    }
}

func toUpperCase(s string) string {
    return strings.ToUpper(s)
}
