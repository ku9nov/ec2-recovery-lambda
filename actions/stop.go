package actions

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
)

type EC2StopInstancesAPI interface {
	StopInstances(ctx context.Context,
		params *ec2.StopInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
}

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func StopInstance(c context.Context, api EC2StopInstancesAPI, input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	resp, err := api.StopInstances(c, input)

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode() == "DryRunOperation" {
		log.Println("User has permission to stop instances.")
		input.DryRun = aws.Bool(false)
		return api.StopInstances(c, input)
	}

	return resp, err
}

func StopInstanceCmd(instanceID *string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.StopInstancesInput{
		InstanceIds: []string{
			*instanceID,
		},
		DryRun: aws.Bool(true),
	}

	_, err = StopInstance(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error stopping the instance")
		log.Println(err)
		return
	}
	parseDur, err := time.ParseDuration("10m")
	if err != nil {
		log.Println("Got an error when parse dur", err)

	}

	test := ec2.NewInstanceStoppedWaiter(client)
	inputDescribe := &ec2.DescribeInstancesInput{InstanceIds: input.InstanceIds}

	out, err := test.WaitForOutput(context.TODO(), inputDescribe, parseDur)
	if err != nil {
		log.Println("Got an error when parse out", err)
	} else {
		_ = out
		log.Println("Successfully stopped instance with ID " + *instanceID)
		StartInstanceCmd(instanceID)
	}

}
