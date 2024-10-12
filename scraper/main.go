package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
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

		// Publish the new properties via SES
		if len(newProperties) > 0 {
			slog.Info("Found new properties!", "count", len(newProperties))
			emailNewProperties(newProperties)
		}
	}

	return "success", nil
}

func emailNewProperties(properties []Property) {
	// Get the AWS region from the environment variable
	region := os.Getenv("AWS_REGION")

	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		slog.Error("Failed to create AWS Session", "error", err)
		return
	}

	// Create an SES client
	svc := ses.New(sess)

	// Email subject and body
	subject := "New properties found!"
	body := buildEmailBody(properties)

	// Construct the email message
	emailAddr := "jonathon_ackers@hotmail.co.uk"
	input := &ses.SendEmailInput{
		Source: aws.String(emailAddr),
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(emailAddr),
			},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(body),
				},
			},
		},
	}

	// Send the email
	result, err := svc.SendEmail(input)
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return
	}

	slog.Info("Successfully sent email!", "message_id", *result.MessageId)
}

func main() {
	lambda.Start(runScraper)
}
