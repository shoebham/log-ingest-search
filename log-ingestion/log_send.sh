#!/bin/bash

# Server URL
SERVER_URL="http://localhost:3000/logs"

# Number of logs to generate
NUM_LOGS=10

for ((i = 0; i < $NUM_LOGS; i++)); do
    # Generate random values for different fields
    LEVEL=("info" "error" "warning")
    RANDOM_LEVEL=${LEVEL[$RANDOM % ${#LEVEL[@]}]}

    MESSAGE="Some random message"

    RESOURCE_ID="server-$((RANDOM % 10000))"

    TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    TRACE_ID=$(cat /dev/urandom | LC_ALL=C tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1)

    SPAN_ID=$(cat /dev/urandom | LC_ALL=C tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)

    COMMIT=$(cat /dev/urandom | LC_ALL=C tr -dc 'a-f0-9' | fold -w 7 | head -n 1)

    PARENT_RESOURCE_ID="server-$((RANDOM % 10000))"

    # Construct JSON log entry
    RANDOM_LOG="{ \"level\": \"$RANDOM_LEVEL\", \"message\": \"$MESSAGE\", \"resourceId\": \"$RESOURCE_ID\", \"timestamp\": \"$TIMESTAMP\", \"traceId\": \"$TRACE_ID\", \"spanId\": \"$SPAN_ID\", \"commit\": \"$COMMIT\", \"metadata\": { \"parentResourceId\": \"$PARENT_RESOURCE_ID\" } }"

    # Send the log entry as a POST request using curl
    curl -X POST -H "Content-Type: application/json" -d "$RANDOM_LOG" "$SERVER_URL"
done
