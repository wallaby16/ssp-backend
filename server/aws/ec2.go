package aws

import (
	"errors"
	"log"
	"net/http"

	"github.com/SchweizerischeBundesbahnen/ssp-backend/server/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gin-gonic/gin"
)

const (
	ec2ListError  = "Instanzen können nicht aufgelistet werden. Bitte erstelle ein Ticket"
	ec2StartError = "Die Instanz konnte nicht gestartet werden. Bitte erstelle ein Ticket"
	ec2StopError  = "Die Instanz konnte nicht gestoppt werden. Bitte erstelle ein Ticket"
)

func listEC2InstancesHandler(c *gin.Context) {
	username := common.GetUserName(c)

	log.Print(username + " lists EC2 Instances")

	instances, err := listEC2InstancesByUsername(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
	} else {
		c.JSON(http.StatusOK, instances)
	}
}

func setEC2InstanceStateHandler(c *gin.Context) {
	username := common.GetUserName(c)
	instanceid := c.Param("instanceid")
	state := c.Param("state")
	log.Print(username + " requested instance " + instanceid + " to " + state)
	instance, err := getInstance(instanceid, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
		return
	}
	account := instance.Account

	switch state {
	case "start":
		res, err := startEC2Instance(instanceid, username, account)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	case "stop":
		res, err := stopEC2Instance(instanceid, username, account)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ApiResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	default:
		c.JSON(http.StatusBadRequest, common.ApiResponse{Message: wrongAPIUsageError})
	}
}

func getInstance(instanceid string, username string) (*common.Instance, error) {
	instances, err := listEC2InstancesByUsername(username)
	if err != nil {
		return nil, err
	}
	for _, instance := range instances.Instances {
		if instance.InstanceId == instanceid {
			return &instance, nil
		}
	}
	log.Println("Could not find an instance with id: " + instanceid)
	return nil, errors.New(ec2ListError)
}

func startEC2Instance(instanceid string, username string, account string) (*common.Instance, error) {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceid),
		},
	}

	var stage string
	if account == accountProd {
		stage = stageProd
	} else {
		stage = stageDev
	}

	svc, err := GetEC2Client(stage)
	if err != nil {
		log.Println("Error getting EC2 client: " + err.Error())
		return nil, errors.New(ec2StartError)
	}

	_, err = svc.StartInstances(input)
	if err != nil {
		log.Println("Error starting EC2 instance (StartInstances API call): " + err.Error())
		return nil, errors.New(ec2StartError)
	}

	filters := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(instanceid),
				},
			},
		},
	}
	err = svc.WaitUntilInstanceRunning(filters)
	if err != nil {
		log.Println("Error waiting for EC2 instance to start: " + err.Error())
		return nil, errors.New(ec2StartError)
	}
	result, err := getInstance(instanceid, username)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func stopEC2Instance(instanceid string, username string, account string) (*common.Instance, error) {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceid),
		},
	}
	var stage string
	if account == accountProd {
		stage = stageProd
	} else {
		stage = stageDev
	}

	svc, err := GetEC2Client(stage)
	if err != nil {
		log.Println("Error getting EC2 client: " + err.Error())
		return nil, errors.New(ec2StopError)
	}

	_, err = svc.StopInstances(input)
	if err != nil {
		log.Println("Error stopping EC2 instance (StopInstances API call): " + err.Error())
		return nil, errors.New(ec2StopError)
	}

	filters := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(instanceid),
				},
			},
		},
	}
	err = svc.WaitUntilInstanceStopped(filters)
	if err != nil {
		log.Println("Error waiting for EC2 instance to stop: " + err.Error())
		return nil, errors.New(ec2StopError)
	}

	result, err := getInstance(instanceid, username)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func listEC2InstancesByUsername(username string) (*common.InstanceListResponse, error) {
	result := common.InstanceListResponse{
		Instances: []common.Instance{},
	}
	nonprodInstances, err := listEC2InstancesByUsernameForAccount(username, accountNonProd)
	if err != nil {
		return nil, err
	}
	prodInstances, err := listEC2InstancesByUsernameForAccount(username, accountProd)
	if err != nil {
		return nil, err
	}

	result.Instances = append(result.Instances, nonprodInstances...)
	result.Instances = append(result.Instances, prodInstances...)

	return &result, nil
}

func listEC2InstancesByUsernameForAccount(username string, account string) ([]common.Instance, error) {
	instances := []common.Instance{}
	filters := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Owner"),
				Values: []*string{
					aws.String("*" + username + "*"),
				},
			},
		},
	}

	var stage string
	if account == accountProd {
		stage = stageProd
	} else {
		stage = stageDev
	}

	svc, err := GetEC2Client(stage)
	if err != nil {
		return nil, errors.New(ec2ListError)
	}

	result, err := svc.DescribeInstances(filters)
	if err != nil {
		log.Print("Unable to list instances (DescribeInstances API call): " + err.Error())
		return nil, errors.New(ec2ListError)
	}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, getInstanceStruct(instance, account))
		}
	}

	return instances, nil
}

func getInstanceStruct(instance *ec2.Instance, account string) common.Instance {
	var name string
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			name = *tag.Value
			break
		}
	}
	return common.Instance{
		Name:       name,
		InstanceId: *instance.InstanceId,
		State:      *instance.State.Name,
		Account:    account,
	}
}
