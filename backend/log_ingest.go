package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/lib/pq"
)

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

func ParseLog(log string) (LogEntry,error){
	var logEntry LogEntry
	err := json.Unmarshal([]byte(log),&logEntry)

	if err!= nil{
		return LogEntry{},err
	}
	return logEntry,nil	
}
var db *sql.DB
var sqliteDB *sql.DB 

func initDbthings(){
	connectDB()
	makeTable()
	connectSQLiteDB()
	makeTableSQLite()
}
func connectSQLiteDB() {
	db, err := sql.Open("sqlite3", "./temp.db")
	if err != nil {
		log.Fatal(err)
	}
	sqliteDB = db
	fmt.Println("Connected to sqlite3")
	
}
func makeTableSQLite(){
	// CREATE TABLE IF NOT EXISTS logs (id SERIAL PRIMARY KEY,level TEXT,message TEXT,resourceId TEXT,timestamp TIMESTAMP,traceId TEXT,spanId TEXT,"commit" TEXT,metadata_parentResourceId TEXT);
	_, err := sqliteDB.Exec(`CREATE TABLE IF NOT EXISTS logs (
		id SERIAL PRIMARY KEY,
		level TEXT,
		message TEXT,
		resourceId TEXT,
		timestamp TIMESTAMP,
		traceId TEXT,
		spanId TEXT,
		"commit" TEXT,
		metadata_parentResourceId TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}
		var count int
	err = sqliteDB.QueryRow("SELECT COUNT(*) FROM logs").Scan(&count)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Number of rows: %d\n", count)
}
func handleLogs(w http.ResponseWriter, r *http.Request){
	fmt.Print("inside handle logs\n")
	var logEntry LogEntry
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err!=nil{
		http.Error(w,"Error parsing log entry",http.StatusBadRequest)
		return
	}


	fmt.Printf("Received log entry: %+v\n", logEntry)

	insertLog(logEntry)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Log received \n"))
}

func writeLogsToFile(logEntry LogEntry){
	fileName := "received_logs.json"
	logs,err := json.MarshalIndent(logEntry,""," ")
	if err!=nil{
		log.Println("Error Marshalling ",err)
		return
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	// Append the log entry to the file
	if _, err := file.WriteString(string(logs) + "\n"); err != nil {
		log.Println("Error writing to file:", err)
	}
}

func greet(w http.ResponseWriter, r *http.Request){	
	fmt.Fprintf(w,"Hi")
	
}


func connectDB(){
	connStr:="user=postgres dbname=temp sslmode=disable"
	dbtemp, err := sql.Open("postgres",connStr)
	
	db=dbtemp
	if err!=nil{
		log.Fatal(err)
	}
	err = db.Ping()

	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Connected to Postgresql")

}
func makeTable(){
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS logs (
		id SERIAL PRIMARY KEY,
		level TEXT,
		message TEXT,
		resourceId TEXT,
		timestamp TIMESTAMP,
		traceId TEXT,
		spanId TEXT,
		commit TEXT,
		metadata_parentResourceId TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}
}

func insertLog(logEntry LogEntry){

	_,err:=db.Exec(`INSERT INTO logs (level,message,resourceId,timestamp,traceId,spanId,commit,metadata_parentResourceID)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
	logEntry.Level,logEntry.Message,logEntry.ResourceID,logEntry.Timestamp,logEntry.TraceID,logEntry.SpanID,logEntry.Commit,logEntry.Metadata.ParentResourceID)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Inserted to postgresql")
	_,err=sqliteDB.Exec(`INSERT INTO logs_fts (level,message,resourceId,timestamp,traceId,spanId,"commit",metadata_parentResourceID)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
	logEntry.Level,logEntry.Message,logEntry.ResourceID,logEntry.Timestamp,logEntry.TraceID,logEntry.SpanID,logEntry.Commit,logEntry.Metadata.ParentResourceID)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Inserted to sqlite")

	
}

type Columns struct {
	Columns []string `json:"columns"`
}

// corsMiddleware is a middleware function to handle CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific headers and methods
		w.Header().Set("Access-Control-Allow-Headers", "hx-target, hx-current-url, hx-trigger, hx-trigger-name, hx-request, hx-prompt, hx-history-restore-request, hx-boosted, Content-Type, Authorization")
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
