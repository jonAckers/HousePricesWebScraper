#!/bin/bash

printf "\n--------------------\n\n"
printf "Building...\n"

# Navigate to the database directory and build
printf "Building the database...\n"
cd ../database || exit
GOOS=linux GOARCH=amd64 go build -o bootstrap
if [ $? -ne 0 ]; then
    printf "Failed to build the database.\n"
    exit 1
fi

# Navigate to the scraper directory and build
printf "Building the scraper...\n"
cd ../scraper || exit
GOOS=linux GOARCH=amd64 go build -o bootstrap
if [ $? -ne 0 ]; then
    printf "Failed to build the scraper.\n"
    exit 1
fi

printf "\n--------------------\n\n"

# Navigate to the cdk directory, download dependencies and run cdk.go
printf "Deploying with CDK...\n"
cd ../cdk || exit
go mod download
go run cdk.go
if [ $? -ne 0 ]; then
    printf "Failed to deploy.\n"
    exit 1
fi

printf "Changes deployed successfully!\n"
