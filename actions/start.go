package actions

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
)

type EC2StartInstancesAPI interface {
	StartInstances(ctx context.Context,
		params *ec2.StartInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
}

func StartInstance(c context.Context, api EC2StartInstancesAPI, input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	resp, err := api.StartInstances(c, input)

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode() == "DryRunOperation" {
		log.Println("User has permission to start an instance.")
		input.DryRun = aws.Bool(false)
		return api.StartInstances(c, input)
	}

	return resp, err
}

func StartInstanceCmd(instanceID *string) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.StartInstancesInput{
		InstanceIds: []string{
			*instanceID,
		},
		DryRun: aws.Bool(true),
	}

	_, err = StartInstance(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error starting the instance")
		log.Println(err)
		return
	}
	SendMessageToSlack(RestartMessage, GreenColor, *instanceID)
	log.Println(RestartMessage, *instanceID)
}
