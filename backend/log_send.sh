#!/bin/bash

# Server URL
SERVER_URL="http://localhost:3000/logs"

# Number of logs to generate
NUM_LOGS=12  # 12 logs for 2 minutes (10 seconds per log)

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

    # Print countdown timer
    for ((remaining = 10; remaining >= 0; remaining--)); do
        echo -ne "Sending next log in $remaining seconds...\033[0K\r"
        sleep 1
    done

    # Send the log entry as a POST request using curl
    curl -X POST -H "Content-Type: application/json" -d "$RANDOM_LOG" "$SERVER_URL"
done
