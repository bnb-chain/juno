package config

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var (
	ErrorNilString = errors.New("secret string is nil")
)

// Params holds all the secret manager information for connecting to database.
// When this information is provided, DBConfiguration.DSN will be replaced with the secret value from AWS SM.
type Params struct {
	SecretId string `yaml:"SecretId"`
	Region   string `yaml:"Region"`
}

func GetString(config *Params) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}

	svc := secretsmanager.New(
		sess,
		aws.NewConfig().WithRegion(config.Region),
	)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(config.SecretId),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	if result.SecretString != nil {
		return *result.SecretString, nil
	}
	return "", ErrorNilString
}
