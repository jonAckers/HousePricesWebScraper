package main

import (
	"database/sql"
	"fmt"
	"log/slog"

	utils "github.com/jonackers/HousePricesWebScraper/utils"
	"github.com/pressly/goose"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

func deployDatabase() error {
	// Get database details from Secrets Manager
	dbDetails, err := utils.GetDbDetails()
	if err != nil {
		slog.Error("Failed to read details from Secrets Manager.")
		return err
	}

	// Create postgres connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
							dbDetails.Host, dbDetails.Port, dbDetails.Username, dbDetails.Password, dbDetails.Name)

	// Open postgres connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Failed to connect to the database.")
		return err
	}
	defer db.Close()

	// Set the Goose migrations directory
	migrationsDir := "migrations"

	// Run Goose migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		slog.Error("Failed to apply migrations.")
		return err
	}

	slog.Info("Database migrations applied successfully!")
	return nil
}

func main() {
	lambda.Start(deployDatabase)
}
