package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func main(){

	file, err := os.Open("template_log.json")
	if err!=nil{
		log.Fatal(err)
	}

	defer file.Close()

	buf := make([]byte, 1024)
	var logs string
	i:=1
	for{
		n,err := file.Read(buf)
		if n>0{
			logs += string(buf[:n])
			i++
		}
		if err!=nil{
			break
		}
	}
	parsedLog,err := ParseLog(logs)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Print(parsedLog.Message)
}
