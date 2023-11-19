
# Log Ingestor and Query Interface


## Backend
Prerequisites: 
```go, postgres, sqlite```

Setup instructions:
1. Clone this project
2. cd to backend directory
3. run ``` go run --tags "fts5" . ``` to run the backend server. 
4. Backend server will now start on port 3000

P.S Make sure you have postgres and sqlite installed, and you have created a db named temp in postgres with the user=postgres. You can change dbname and user in [connectDb()](https://github.com/shoebham/log-ingest-search/blob/main/backend/db.go#L26) in db.go. 
For sqlite a db with the name temp.db is created when the server is run.

## Frontend
for the query interface just open the index.html and your frontend will be ready.

## Demo
### Server start
<img width="1552" alt="image" src="https://github.com/shoebham/log-ingest-search/assets/25881429/5bc16be4-a6fe-497c-8044-624ae53986ef">

### Sending a log
<img width="1392" alt="image" src="https://github.com/shoebham/log-ingest-search/assets/25881429/1ea73e1c-8cd4-49b8-96e1-b9836bf0eaa0">
<img width="1005" alt="image" src="https://github.com/shoebham/log-ingest-search/assets/25881429/0c397a44-0a7b-4b43-9482-f56d0e07802f">


### Query Interface
<img width="1552" alt="image" src="https://github.com/shoebham/log-ingest-search/assets/25881429/62521bbe-fc2b-4d2b-8d18-4f49bedd7332">

### Searching a log
<img width="1552" alt="image" src="https://github.com/shoebham/log-ingest-search/assets/25881429/9e6edd93-089b-44d4-a71b-0dd0dd938f0b">

### Filter search
![filter-search](https://github.com/shoebham/log-ingest-search/assets/25881429/910b2184-865b-4779-8b22-8d29a14bae1b)



### full text search
![full-text-search](https://github.com/shoebham/log-ingest-search/assets/25881429/4df1a3f0-3c9c-42a4-ae33-d0400392d583)




