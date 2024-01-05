#!/bin/bash

set -e

ENV_FILE=".env"
BINARY_NAME="app"

if [ ! -f "$ENV_FILE" ]; then
    echo "Missing .env file at $(pwd)/$ENV_FILE"
    exit 1
fi

export $(grep -v '^#' $ENV_FILE | xargs)

echo "Building the Go application..."
go build -o $BINARY_NAME cmd/main.go

echo "Running the application..."
./$BINARY_NAME