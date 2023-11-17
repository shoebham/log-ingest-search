package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type SearchCriteria struct {
	Criteria []interface{} `json:"criteria"`
}


type SearchParam struct {

	Column  string `json:"column"`
	Operand string `json:"operand"`
	Value   string `json:"value"`
}

func search(w http.ResponseWriter, r *http.Request){
	var searchCriteria SearchCriteria

	err := json.NewDecoder(r.Body).Decode(&searchCriteria)
	if err != nil {
		http.Error(w, "Error parsing search criteria", http.StatusBadRequest)
		return
	}

	query := constructQuery(searchCriteria)
	fmt.Println("QUERY:", query)


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