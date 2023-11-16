package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

func handleLogs(w http.ResponseWriter, r *http.Request){
	fmt.Print("inside handle logs\n")
	var logEntry LogEntry
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err!=nil{
		http.Error(w,"Error parsing log entry",http.StatusBadRequest)
		return
	}


	fmt.Printf("Received log entry: %+v\n", logEntry)
	writeLogsToFile(logEntry)
	// Respond with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Log received and written to file successfully"))
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


func main(){

	http.HandleFunc("/",greet)
	http.HandleFunc("/logs",handleLogs)

	fmt.Print("Listenting on port 3000")

	err := http.ListenAndServe(":3000",nil)
	if err!=nil{
		log.Fatal("Server Error",err)
	}
}
