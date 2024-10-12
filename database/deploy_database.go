package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type SecretToken struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func deployDatabase() error {
	// Get database details from Secrets Manager
	dbDetails, err := getDbDetails()
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

func getDbDetails() (SecretToken, error) {
	// Load secret name from environment variables
	secretName := os.Getenv("SECRET_NAME")
	curRegion := os.Getenv("AWS_REGION")

	// Load the default AWS SDK config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(curRegion))
	if err != nil {
		slog.Error("Failed to load AWS SDK config.")
		return SecretToken{}, err
	}

	// Create a Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	// Get the secret value from Secrets Manager
	secretValue, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		slog.Error("Failed to get database secret.")
		return SecretToken{}, err
	}

	// Read the password from the secret
	var result SecretToken
	if err := json.Unmarshal([]byte(*secretValue.SecretString), &result); err != nil {
		slog.Error("Failed to parse details from secret.")
		return SecretToken{}, err
	}

	// Return details
	return result, nil
}

func main() {
	lambda.Start(deployDatabase)
}
