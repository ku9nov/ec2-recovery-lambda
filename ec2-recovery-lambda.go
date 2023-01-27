package main

import (
	"context"
	"log"

	"ec2-recovery-lambda/actions"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2DescribeInstanceStatusAPI interface {
	DescribeInstanceStatus(ctx context.Context,
		params *ec2.DescribeInstanceStatusInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error)
}

func GetInstances(c context.Context, api EC2DescribeInstanceStatusAPI, input *ec2.DescribeInstanceStatusInput) (*ec2.DescribeInstanceStatusOutput, error) {
	return api.DescribeInstanceStatus(c, input)
}

func getInstanceName(instanceID string, cfg aws.Config) (instanceNameById string) {

	ec2Svc := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := ec2Svc.DescribeInstances(context.Background(), input)
	if err != nil {
		log.Println("Error", err)
		return
	}

	// Get the instance name
	instanceName := ""
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					instanceName = *tag.Value
					break
				}
			}
		}
	}
	log.Println("Instance Name:", instanceName)
	return instanceName
}

func checkInstances() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstanceStatusInput{}

	result, err := GetInstances(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error retrieving information about your Amazon EC2 instances:")
		log.Println(err)
		return
	}
	instanceCount := len(result.InstanceStatuses)
	log.Println("Count of running instances: ", instanceCount)
	for _, r := range result.InstanceStatuses {
		if r.InstanceStatus.Status == "impaired" || r.SystemStatus.Status == "impaired" {
			log.Println(actions.AlarmMessage, *r.InstanceId)
			nameOfImpairedInstance := getInstanceName(*r.InstanceId, cfg)
			actions.SendMessageToSlack(actions.AlarmMessage, actions.RedColor, nameOfImpairedInstance, *r.InstanceId, instanceCount)
			actions.StopInstanceCmd(r.InstanceId, instanceCount, nameOfImpairedInstance)
		}
	}
}

func main() {
	lambda.Start(checkInstances)
}
