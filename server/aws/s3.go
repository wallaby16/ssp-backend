package aws

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/oscp/cloud-selfservice-portal/server/common"
)

const (
	genericAPIError    = "Fehler beim Aufruf der S3-API. Bitte erstelle ein Ticket"
	genericListError   = "Ressourcen können nicht aufgelistet werden. Bitte erstelle ein Ticket"
	genericCreateError = "Erstellung des Buckets fehlgeschlagen. Bitte erstelle ein Ticket"
	wrongAPIUsageError = "Invalid API call - parameters did not match to method definition"
)

// RegisterRoutes registers the routes for S3
func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/aws/s3/list", listS3BucketsHandler)
	r.POST("/aws/s3/new", newS3BucketHandler)
	r.POST("/aws/s3/newreaduser", newS3ReadUserHandler)
	r.POST("/aws/s3/newwrite", newS3WriteUserHandler)
}

func newS3BucketHandler(c *gin.Context) {
	username := common.GetUserName(c)

	var data common.NewS3BucketCommand
	if c.BindJSON(&data) == nil {

		newbucketname := generateS3Bucketname(data.BucketName, data.Stage)
		if err := validateNewS3Bucket(username, data.Project, newbucketname, data.Billing, data.Stage); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		log.Print("Creating new bucket " + newbucketname + " for " + username)
		if err := createNewS3Bucket(username, data.Project, newbucketname, data.Billing, data.Stage); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, common.ApiResponse{
				Message: "Es wurde ein neuer S3 Bucket erstellt: " + newbucketname,
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
	}
}

func listS3BucketsHandler(c *gin.Context) {
	username := common.GetUserName(c)

	log.Print(username + " lists S3 buckets")
	myBuckets, _ := listS3BucketByUsername(username)

	c.JSON(http.StatusOK, myBuckets)
}

func newS3ReadUserHandler(c *gin.Context) {
	username := common.GetUserName(c)

	var data common.NewS3UserCommand
	if c.BindJSON(&data) == nil {

		if err := validateNewS3User(username, data.BucketName, data.UserName); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		log.Print(username + " creates a new read user (" + data.UserName + ") for " + data.BucketName)
		credentials, err := createNewS3ReadUser(data.BucketName, data.UserName)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, credentials)
		}
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
	}
}

func newS3WriteUserHandler(c *gin.Context) {
	username := common.GetUserName(c)

	var data common.NewS3UserCommand
	if c.BindJSON(&data) == nil {

		if err := validateNewS3User(username, data.BucketName, data.UserName); err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}

		log.Print(username + " creates a new read/write user (" + data.UserName + ") for " + data.BucketName)
		credentials, err := createNewS3WriteUser(data.BucketName, data.UserName)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusOK, credentials)
		}
	} else {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
	}
}

func validateNewS3Bucket(username string, projectname string, bucketname string, billing string, stage string) error {

	if len(stage) == 0 {
		return errors.New("Umgebung muss definiert werden")
	}
	if len(billing) == 0 {
		return errors.New("Verrechnungsnummer muss definiert sein")
	}
	if len(bucketname) == 0 {
		return errors.New("Bucketname muss definiert sein")
	}
	if len(projectname) == 0 {
		return errors.New("Projekt muss definiert sein")
	}

	if len(bucketname) > 63 {
		// http://docs.aws.amazon.com/AmazonS3/latest/dev/BucketRestrictions.html
		return errors.New("Generierter Bucketname " + bucketname + " ist zu lang")
	}
	var validName = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`).MatchString
	if !validName(bucketname) {
		return errors.New("Bucketname kann nur alphanumerische Zeichen und Bindestriche enthalten")
	}

	if stage != "dev" && stage != "test" && stage != "int" && stage != "prod" {
		return errors.New("Unbekannte Umgebung: " + stage)
	}

	// Check if bucket already exists
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_S3_REGION"))},
	)
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		log.Print("Error while trying to validate new bucket (ListBucket call): " + err.Error())
		return errors.New(genericCreateError)
	}

	for _, b := range result.Buckets {
		if *b.Name == bucketname {
			log.Print("Error, bucket " + bucketname + "already exists")
			return errors.New("Fehler: Bucket " + bucketname + " existiert bereits")
		}
	}

	// Everything OK
	return nil
}

func createNewS3Bucket(username string, projectname string, bucketname string, billing string, stage string) error {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_S3_REGION"))},
	)

	// Create S3 service client
	svc := s3.New(sess)

	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketname),
	})
	if err != nil {
		log.Print("Error on CreateBucket call (username=" + username + ", bucketname=" + bucketname + "): " + err.Error())
		return errors.New(genericCreateError)
	}

	// Wait until bucket is created before finishing
	log.Print("Waiting for bucket " + bucketname + " to be created...")
	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketname),
	})

	if err != nil {
		log.Print("Error when creating S3 bucket in WaitUntilBucketExists: " + err.Error())
		return errors.New(genericCreateError)
	}

	_, err = svc.PutBucketTagging(&s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketname),
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
				{Key: aws.String("Creator"), Value: aws.String(username)},
				{Key: aws.String("Project"), Value: aws.String(projectname)},
				{Key: aws.String("Accounting_Number"), Value: aws.String(billing)},
				{Key: aws.String("Stage"), Value: aws.String(stage)},
			},
		}})
	if err != nil {
		log.Print("Tagging bucket " + bucketname + " failed: " + err.Error())
		return errors.New(genericCreateError)
	}

	//log.Print("Creating IAM policies for bucket " + bucketname + "...")
	// Create a IAM service client.
	iamSvc := iam.New(sess)

	readPolicy := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			StatementEntry{
				Effect: "Allow",
				Action: []string{
					"s3:Get*",  // Allow Get commands
					"s3:List*", // Allow List commands
				},
				Resource: "arn:aws:s3:::" + bucketname,
			},
			StatementEntry{
				Effect: "Allow",
				Action: []string{
					"s3:Get*",  // Allow Get commands
					"s3:List*", // Allow List commands
				},
				Resource: "arn:aws:s3:::" + bucketname + "/*",
			},
		},
	}

	writePolicy := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			StatementEntry{
				Effect: "Allow",
				Action: []string{
					"s3:Get*",    // Allow Get commands
					"s3:List*",   // Allow List commands
					"s3:Put*",    // Allow Put commands
					"s3:Delete*", // Allow Delete commands
				},
				Resource: "arn:aws:s3:::" + bucketname,
			},
			StatementEntry{
				Effect: "Allow",
				Action: []string{
					"s3:Get*",    // Allow Get commands
					"s3:List*",   // Allow List commands
					"s3:Put*",    // Allow Put commands
					"s3:Delete*", // Allow Delete commands
				},
				Resource: "arn:aws:s3:::" + bucketname + "/*",
			},
		},
	}

	// Read policy
	b, err := json.Marshal(&readPolicy)
	if err != nil {
		log.Print("Error marshaling readPolicy: " + err.Error())
		return errors.New(genericCreateError)
	}

	_, err = iamSvc.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(b)),
		PolicyName:     aws.String(bucketname + "-BucketReadPolicy"),
	})
	if err != nil {
		log.Print("Error CreatePolicy for BucketReadPolicy failed: " + err.Error())
		return errors.New(genericCreateError)
	}

	// Write policy
	c, err := json.Marshal(&writePolicy)
	if err != nil {
		log.Print("Error marshaling writePolicy: " + err.Error())
		return errors.New(genericCreateError)
	}

	_, err = iamSvc.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(c)),
		PolicyName:     aws.String(bucketname + "-BucketWritePolicy"),
	})
	if err != nil {
		log.Print("Error CreatePolicy for BucketWritePolicy failed: " + err.Error())
		return errors.New(genericCreateError)
	}

	log.Print("Bucket " + bucketname + " and IAM policies successfully created")
	return nil
}

func generateS3Bucketname(bucketname string, stage string) string {
	// Generate bucketname: <prefix>-<bucketname>-<stage_suffix>
	bucketPrefix := os.Getenv("AWS_S3_BUCKET_PREFIX")
	stageSuffix := "nonprod"
	if stage == "prod" {
		stageSuffix = "prod"
	}
	return strings.ToLower(bucketPrefix + bucketname + "-" + stageSuffix)
}

func listS3BucketByUsername(username string) (common.BucketListResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_S3_REGION"))},
	)

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		log.Print("Unable to list buckets (ListBuckets API call): " + err.Error())
		return common.BucketListResponse{}, errors.New(genericListError)
	}

	// Return list
	responseList := common.BucketListResponse{
		Buckets: []string{},
	}

	for _, b := range result.Buckets {
		// Get bucket tags
		taggingParams := &s3.GetBucketTaggingInput{
			Bucket: aws.String(*b.Name),
		}
		result, err := svc.GetBucketTagging(taggingParams)
		if err != nil {
			log.Print("Unable to get tags for bucket, username " + username + ": " + err.Error())
			return common.BucketListResponse{}, errors.New(genericListError)
		}

		// Get list of buckets where "Creator" equals username and return only those
		for _, tag := range result.TagSet {
			if *tag.Key == "Creator" && *tag.Value == username {
				responseList.Buckets = append(responseList.Buckets, *b.Name)
			}
		}
	}
	return responseList, nil
}
