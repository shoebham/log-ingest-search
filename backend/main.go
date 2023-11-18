package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/lib/pq"
)

type Columns struct {
	Columns []string `json:"columns"`
}
type LogEntry struct {
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	ResourceID  string    `json:"resourceId"`
	Timestamp   time.Time `json:"timestamp"`
	TraceID     string    `json:"traceId"`
	SpanID      string    `json:"spanId"`
	Commit      string    `json:"commit"`
	Metadata    struct {
		ParentResourceID string `json:"parentResourceId"`
	} `json:"metadata"`
}

var db *sql.DB
var sqliteDB *sql.DB 



func greet(w http.ResponseWriter, r *http.Request){	
	fmt.Fprintf(w,"Hi")
	
}


func ParseLog(log string) (LogEntry,error){
	var logEntry LogEntry
	err := json.Unmarshal([]byte(log),&logEntry)

	if err!= nil{
		return LogEntry{},err
	}
	return logEntry,nil	
}


// corsMiddleware is a middleware function to handle CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific headers and methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Allow credentials if needed (set to "true" to allow)
		w.Header().Set("Access-Control-Allow-Credentials", "false")

		// Preflight OPTIONS request handling
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main(){
	initDbthings()
	http.HandleFunc("/",greet)
	http.HandleFunc("/logs",handleLogs)
	http.HandleFunc("/columns",fetchColumnsHandler)
	http.HandleFunc("/search",search)
	http.HandleFunc("/searchRealTime", debounceAPIRequest(realTimeSearch))
	fmt.Println("Listenting on port 3000")
	handler := corsMiddleware(http.DefaultServeMux)

	// Start server with the CORS-wrapped handler
	server := &http.Server{
		Addr:           ":3000",
		Handler:        handler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server Error: ", err)
	}
	defer db.Close()
	defer sqliteDB.Close()

}
