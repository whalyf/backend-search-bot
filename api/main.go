package main
// package handler
//package handler  TO VERCEL

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
    "github.com/resend/resend-go/v2"
		g "github.com/serpapi/google-search-results-golang"
)

func main() {
    router := mux.NewRouter()

		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

    // Register the handler for the "/process" endpoint with Gorilla Mux
    router.HandleFunc("/process", HandleProcessRequest).Methods(http.MethodPost, http.MethodOptions)
    router.HandleFunc("/", Greetings).Methods("GET")

    corsHandler := handlers.CORS(
        handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With", "Authorization"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        handlers.AllowedOrigins([]string{"*"}), // Allow requests from any origin
    )

    // Wrap your router with the CORS middleware
    http.Handle("/process", corsHandler(router))
    http.Handle("/", corsHandler(router))

    log.Fatal(http.ListenAndServe(":5555", corsHandler(router)))
}

func Greetings(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<div><h1>Welcome to Google Digger GolangApi</h1><span>Código Fonte: <a target='_blank' href='https://github.com/whalyf/backend-search-bot'>Aqui!</a></span></div>")
}

func HandleProcessRequest(w http.ResponseWriter, r *http.Request) {
    var searchParams map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&searchParams)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    processedResults := processNestJSData(searchParams)
		keywords, ok := processedResults["keywords"].(string)
		if !ok {
				http.Error(w, "Keywords not found", http.StatusInternalServerError)
				return
		}

		email, ok := processedResults["email"].(string)
		if !ok {
				http.Error(w, "Email not found", http.StatusInternalServerError)
				return
		}
    // BUSCA É INVOCADA UTILIZANDO SERP_API
    searchResult := searchOnGoogle(keywords)
    // JSON DE RESPOSTA É CONVERTIDO EM HTML
    htmlFormat := prettyPrintHTML(searchResult)

    // EMAIL ENVIADO COM O HTMLJSON DA BUSCA
    sendEmail(email, htmlFormat)

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173/")
    json.NewEncoder(w).Encode(searchResult)
}

func sendEmail(email string, prettyJSON string) {
  client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

  params := &resend.SendEmailRequest{
      From:    "Google Search Digger <google-digger@resend.dev>",
      To:      []string{email},
      Html:    fmt.Sprintf("<div>%s</div>", prettyJSON),
      Subject: "Resultado das buscas",
  }

  _, err := client.Emails.Send(params)
  if err != nil {
      fmt.Println(err.Error())
      return
  }
}

func prettyPrintHTML(data interface{}) string {
  // Convert data to JSON string
  jsonData, err := json.MarshalIndent(data, "", "  ")
  if err != nil {
      fmt.Fprintf(os.Stderr, "Error JSON: %v\n", err)
      return ""
  }

  prettyJSON := strings.ReplaceAll(string(jsonData), " ", "&nbsp;")
  prettyJSON = strings.ReplaceAll(prettyJSON, "\n", "<br/>")

  prettyJSON = strings.ReplaceAll(prettyJSON, "<br/>{", "<br/>&nbsp;&nbsp;{")
  prettyJSON = strings.ReplaceAll(prettyJSON, "<br/>&nbsp;&nbsp;}", "<br/>}")

  return prettyJSON
}

func searchOnGoogle(keywords string) map[string]interface{}{
	parameter := map[string]string{
    "api_key": os.Getenv("SERPAPI_KEY"),
    "engine": "google",
    "q": keywords,
    "location": "Brazil",
    "google_domain": "google.com.br",
    "gl": "br",
    "hl": "pt",
  }
  search := g.NewGoogleSearch(parameter, os.Getenv("SERPAPI_KEY"))
  results, err := search.GetJSON()

	fmt.Println(err)
	return results
}

func processNestJSData(searchParams map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        // "location":  searchParams["location"],
        // "frequency": searchParams["frequency"],
        "email": searchParams["email"],
        "keywords":  toUpperCase(searchParams["keywords"].(string)),
        "dateTime":  time.Now().Format(time.RFC3339),
        "searchId":   searchParams["searchId"],
    }
}

func toUpperCase(s string) string {
    return strings.ToUpper(s)
}
