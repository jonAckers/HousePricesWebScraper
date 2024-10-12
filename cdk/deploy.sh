#!/bin/bash

echo "\n--------------------\n"
echo "Building..."

# Navigate to the database directory and build
echo "Building the database..."
cd ../database || exit
GOOS=linux GOARCH=amd64 go build -o bootstrap
if [ $? -ne 0 ]; then
    echo "Failed to build the database."
    exit 1
fi

# Navigate to the scraper directory and build
echo "Building the scraper..."
cd ../scraper || exit
GOOS=linux GOARCH=amd64 go build -o bootstrap
if [ $? -ne 0 ]; then
    echo "Failed to build the scraper."
    exit 1
fi

echo "\n--------------------\n"

# Navigate to the cdk directory, download dependencies and run cdk.go
echo "Deploying with CDK..."
cd ../cdk || exit
go mod download
go run cdk.go
if [ $? -ne 0 ]; then
    echo "Failed to deploy."
    exit 1
fi

echo "Changes deployed successfully!"
