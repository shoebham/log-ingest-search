package main

import (
	"testing"
)


func BenchmarkLogParsing(b *testing.B){

	logData := `{
	"level": "error",
	"message": "Failed to connect to DB",
    "resourceId": "server-1234",
	"timestamp": "2023-09-15T08:00:00Z",
	"traceId": "abc-xyz-123",
    "spanId": "span-456",
    "commit": "5e5342f",
    "metadata": {
        "parentResourceId": "server-0987"
    }
}`

	for i:=0;i<b.N;i++{
		_,err := ParseLog(logData)
		if err != nil{
			b.Errorf("Error parsing log entry: %s",err)
		}
	}
}