package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func initDbthings(){
	connectDB()
	makeTable()
	connectSQLiteDB()
	makeTableSQLite()
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


