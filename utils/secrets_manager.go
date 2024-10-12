package utils

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
)

type SecretToken struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetDbDetails() (SecretToken, error) {
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
