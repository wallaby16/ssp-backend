package aws

import (
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/oscp/cloud-selfservice-portal/server/common"
)

const (
	genericUserCreationError = "Es ist ein Fehler bei der Erstellung des Benutzers aufgetreten"
)

// PolicyDocument IAM Policy Document
type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

// StatementEntry IAM Statement Entry
type StatementEntry struct {
	Effect   string
	Action   []string
	Resource string
}

func validateNewS3User(username string, bucketname string, newuser string) error {

	if len(username) == 0 {
		return errors.New("Username must be specified")
	}
	if len(bucketname) == 0 {
		return errors.New("Bucket name must be specified")
	}
	if len(newuser) == 0 {
		return errors.New("Name of new user must be specified")
	}

	if (len(newuser) + len(bucketname)) > 63 {
		// http://docs.aws.amazon.com/IAM/latest/UserGuide/reference_iam-limits.html
		return errors.New("Generierter Benutzername '" + bucketname + "-" + newuser + "' ist zu lang")
	}
	var validName = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString
	if !validName(bucketname) {
		return errors.New("Benutzername kann nur alphanumerische Zeichen und Bindestriche enthalten")
	}

	// Check if user already exists
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_S3_REGION"))},
	)
	svc := iam.New(sess)
	result, err := svc.ListUsers(nil)
	if err != nil {
		log.Print("Error while trying to create a new user (ListUsers call): " + err.Error())
		return errors.New(genericUserCreationError)
	}
	// Loop over existing users
	for _, u := range result.Users {
		if *u.UserName == newuser {
			log.Print("Error, user " + newuser + "already exists")
			return errors.New("Fehler: IAM-Benutzer " + newuser + " existiert bereits")
		}
	}

	// Make sure the user is allowed to create new IAM users for this bucket
	myBuckets, _ := listS3BucketByUsername(username)
	for _, mybucketname := range myBuckets.Buckets {
		if bucketname == mybucketname {
			// Everything OK
			return nil
		}
	}
	return errors.New("User " + username + " does not own bucket " + bucketname)
}

func createNewS3ReadUser(bucketname string, s3username string) (common.S3CredentialsResponse, error) {
	generatedName := bucketname + "-" + s3username
	cred, err := createNewS3User(generatedName)
	if err != nil {
		log.Print("Error while calling createNewS3User: " + err.Error())
		return cred, errors.New(genericUserCreationError)
	}

	err = attachIAMPolicyToUser(bucketname+"-BucketReadPolicy", generatedName)
	if err != nil {
		log.Print("Error while calling attachIAMPolicyToUser: " + err.Error())
		return cred, errors.New(genericUserCreationError)
	}
	return cred, nil
}

func createNewS3WriteUser(bucketname string, s3username string) (common.S3CredentialsResponse, error) {
	generatedName := bucketname + "-" + s3username
	cred, err := createNewS3User(generatedName)
	if err != nil {
		log.Print("Error while calling createNewS3User: " + err.Error())
		return cred, errors.New(genericUserCreationError)
	}

	err = attachIAMPolicyToUser(bucketname+"-BucketWritePolicy", generatedName)
	if err != nil {
		log.Print("Error while calling attachIAMPolicyToUser: " + err.Error())
		return cred, errors.New(genericUserCreationError)
	}
	return cred, nil
}

func createNewS3User(name string) (common.S3CredentialsResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_S3_REGION"))},
	)

	// Create a IAM service client.
	svc := iam.New(sess)

	_, err = svc.GetUser(&iam.GetUserInput{
		UserName: aws.String(name),
	})

	var cred common.S3CredentialsResponse
	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == iam.ErrCodeNoSuchEntityException {
		_, err := svc.CreateUser(&iam.CreateUserInput{
			UserName: aws.String(name),
		})

		if err != nil {
			return cred, errors.New("CreateUser error in createNewS3User: " + err.Error())
		}

		// Create access key
		result, err := svc.CreateAccessKey(&iam.CreateAccessKeyInput{
			UserName: aws.String(name),
		})
		cred.AccessKeyID = *result.AccessKey.AccessKeyId
		cred.SecretKey = *result.AccessKey.SecretAccessKey
		return cred, nil
	}
	return cred, errors.New("GetUser error: " + err.Error())
}

func attachIAMPolicyToUser(policyName string, username string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_S3_REGION"))},
	)

	// Create a IAM service client.
	svc := iam.New(sess)

	// First, figure out the AWS account ID
	result, err := svc.GetUser(nil)
	var accountNumber string
	if err != nil {
		return errors.New("GetUser error in attachIAMPolicyToUser() while trying to determine account ID: " + err.Error())
	}
	re := regexp.MustCompile("[0-9]+")
	accountNumber = re.FindString(*result.User.Arn)

	// Then, attach the policy given to the user
	input := &iam.AttachUserPolicyInput{
		PolicyArn: aws.String("arn:aws:iam::" + accountNumber + ":policy/" + policyName),
		UserName:  aws.String(username),
	}
	_, err = svc.AttachUserPolicy(input)
	if err != nil {
		return errors.New("AttachUserPolicy error: " + err.Error())
	}
	return nil
}
