package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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


func search(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query()
	var queryParams []string
	var args []interface{}

	for key, values := range query {
		if key != "" && len(values) > 0 {
			for _, value := range values {
				queryParams = append(queryParams, fmt.Sprintf("%s=$%d", key, len(args)+1))
				args = append(args, value)
			}
		}
	}

	whereClause := strings.Join(queryParams, " OR ")
	sqlQuery := fmt.Sprintf("SELECT * FROM logs WHERE %s;", whereClause)

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Println(err)
			continue
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}
		results = append(results, rowMap)
	}

	// Convert results to JSON
	responseJSON, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

