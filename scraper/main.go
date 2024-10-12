package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jonackers/HousePricesWebScraper/scraper/internal/database"
	"github.com/jonackers/HousePricesWebScraper/utils"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	db *database.Queries
	ctx context.Context
}

func runScraper(ctx context.Context) (string, error) {
	// Fetch database secret from Secrets Manager
	dbDetails, err := utils.GetDbDetails()
	if err != nil {
		slog.Error("Failed to read details from Secrets Manager.")
		return "Secrets Manager failure", err
	}

	// Create postgres connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
							dbDetails.Host, dbDetails.Port, dbDetails.Username, dbDetails.Password, dbDetails.Name)

	// Open postgres connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Failed to connect to the database.")
		return "Database connection failure", err
	}
	defer db.Close()

	// Create database connection
	dbQueries := database.New(db)
	cfg := &dbConfig {
		db: dbQueries,
		ctx: ctx,
	}

	// Get properties from Rightmove
	foundProperties := parseHouseDetails()
	if len(foundProperties) > 0 {
		newProperties, err := cfg.getNewProperties(foundProperties)
		if err != nil {
			return "Scraping failure", err
		}

		slog.Info("Found new properties!", "count", len(newProperties))
	}

	return "success", nil
}

func main() {
	lambda.Start(runScraper)
}
