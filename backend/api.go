package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)
type SearchCriteria struct {
	Criteria []interface{} `json:"criteria"`
}


type SearchParam struct {

	Column  string `json:"column"`
	Operand string `json:"operand"`
	Value   string `json:"value"`
}

// comes from /columns endpoint
func fetchColumnsHandler(w http.ResponseWriter, r *http.Request) {
	
	// Query to fetch column names from a specific table
	rows, err := db.Query("SELECT column_name FROM information_schema.columns WHERE table_name = 'logs'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var columns []string

	// Iterate through rows and collect column names
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			log.Fatal(err)
		}
		columns = append(columns, columnName)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	columnData := Columns{Columns: columns}

	// Convert columns to JSON
	columnsJSON, err := json.Marshal(columnData)
	if err != nil {
		log.Fatal(err)
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(columnsJSON)
}

// comes from /logs endpoint
func handleLogs(w http.ResponseWriter, r *http.Request){
	fmt.Println("Log receieved at",time.Now().Local())
	var logEntry LogEntry
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err!=nil{
		http.Error(w,"Error parsing log entry",http.StatusBadRequest)
		return
	}
	// fmt.Printf("Received log entry: %+v\n", logEntry)
	insertLog(logEntry)
	fmt.Println()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Log received \n"))
}



var (
	lastRequest   time.Time
	lastRequestMu sync.Mutex
)

const debounceDuration = 500 * time.Millisecond // Define your desired debounce duration

// debounceAPIRequest ensures that the API endpoint is not called more frequently than the debounce duration
func debounceAPIRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastRequestMu.Lock()
		defer lastRequestMu.Unlock()

		// Calculate the duration since the last request
		elapsed := time.Since(lastRequest)

		// If the duration is less than the debounce duration, wait
		if elapsed < debounceDuration {
			time.Sleep(debounceDuration - elapsed)
		}

		lastRequest = time.Now()

		// Process the request
		next(w, r)
	})
}
func realTimeSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Empty query", http.StatusBadRequest)
		return
	}

	// Perform real-time search logic using the 'query'...
	
	// Execute the search query
	rows, err := sqliteDB.Query("SELECT rowid,* FROM logs_fts WHERE logs_fts MATCH ?", query+"*")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Process query results
	results := processQueryResults(rows)
	
	// Respond with the search results
	responseData := map[string]interface{}{
		"count":  len(results),
		"result": results,
	}

	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
func search(w http.ResponseWriter, r *http.Request){
	var searchCriteria SearchCriteria

	err := json.NewDecoder(r.Body).Decode(&searchCriteria)
	if err != nil {
		http.Error(w, "Error parsing search criteria", http.StatusBadRequest)
		return
	}

	query := constructQuery(searchCriteria)



	rows, err := db.Query(query)
	fmt.Println("QUERY: ",query )
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := processQueryResults(rows)

	// Include the count in the response data
	responseData := map[string]interface{}{
		"count":  len(results),
		"result": results,
	}
	// Respond with the search results
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}


func constructQuery(searchCriteria SearchCriteria) string {
	query := "SELECT * FROM logs WHERE "
	// fmt.Printf("%+v",searchCriteria)
	if len(searchCriteria.Criteria)==0{
		query+="true";
	}
	for i, criteria := range searchCriteria.Criteria {
		if i > 0 {
			// Check if the criteria is a logical operator (AND/OR)
			if logical, ok := criteria.(string); ok {
				query += fmt.Sprintf(" %s ", logical)
				continue
			}
		}

		param, ok := criteria.(map[string]interface{})
		if !ok {
			log.Println("Invalid search parameter format")
			continue
		}

		column, _ := param["column"].(string)
		operand, _ := param["operand"].(string)
		value, _ := param["value"].(string)
		// Check for regex operand
		if operand == "=~" {
			query += fmt.Sprintf("%s ~ '%s'", column, value)
		} else {
			query += fmt.Sprintf("%s %s '%s'", column, operand, value)
		}
		// query += fmt.Sprintf("%s %s '%s'", column, operand, value)
	}
	return query
}

func processQueryResults(rows *sql.Rows) []map[string]interface{} {
	var results []map[string]interface{}

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePointers := make([]interface{}, count)
	for i := range columns {
		valuePointers[i] = &values[i]
	}

	for rows.Next() {
		rows.Scan(valuePointers...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}
		results = append(results, entry)
	}
	return results
}