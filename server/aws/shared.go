package aws

import (
	"os"

	"log"

	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	wrongAPIUsageError = "Invalid API call - parameters did not match to method definition"
	genericAwsAPIError = "Fehler beim Aufruf der AWS API. Bitte erstelle ein Ticket"
)

const (
	accountProd    = "prod"
	accountNonProd = "nonprod"
)

const (
	stageDev  = "dev"
	stageTest = "test"
	stageInt  = "int"
	stageProd = "prod"
)

const (
	bucketReadPolicy = "-BucketReadPolicy"
	bucketWritePolicy = "-BucketWritePolicy"
)

func GetS3Client(stage string) (*s3.S3, error) {
	account, err := getAccountForStage(stage)
	if err != nil {
		return nil, err
	}

	sess, err := getAwsSession(account)
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}

func GetIAMClient(stage string) (*iam.IAM, error) {
	account, err := getAccountForStage(stage)
	if err != nil {
		return nil, err
	}

	sess, err := getAwsSession(account)
	if err != nil {
		return nil, err
	}
	return iam.New(sess), nil
}

func getAwsSession(account string) (*session.Session, error) {
	// Validate necessary env variables
	region := os.Getenv("AWS_REGION")
	if len(region) == 0 {
		log.Fatal("Env variable 'AWS_REGION' must be specified")
	}
	bucketPrefix := os.Getenv("AWS_S3_BUCKET_PREFIX")
	if len(bucketPrefix) == 0 {
		log.Fatal("Env variable 'AWS_S3_BUCKET_PREFIX' must be specified")
	}

	// Create AWS session based on account
	var accessKeyID string
	var accessSecret string

	switch account {
	case accountProd:
		accessKeyID = os.Getenv("AWS_PROD_ACCESS_KEY_ID")
		accessSecret = os.Getenv("AWS_PROD_SECRET_ACCESS_KEY")
	case accountNonProd:
		accessKeyID = os.Getenv("AWS_NONPROD_ACCESS_KEY_ID")
		accessSecret = os.Getenv("AWS_NONPROD_SECRET_ACCESS_KEY")
	default:
		log.Println("Invalid account: " + account)
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, accessSecret, ""),
		Region:      aws.String(region)},
	)

	if err != nil {
		log.Println("Error creating aws session: ", err.Error())
		return nil, errors.New(genericAwsAPIError)
	}

	return sess, nil
}

// getAccountForStage remapps the stage string form the UI to
// the technical AWS account
// dev, test, int = NONPROD
// prod = PROD
func getAccountForStage(stage string) (string, error) {
	switch stage {
	case stageDev, stageTest, stageInt:
		return accountNonProd, nil
	case stageProd:
		return accountProd, nil
	default:
		log.Println("Could not map to account, invalid stage: " + stage)
		return "", errors.New(wrongAPIUsageError)
	}
}
