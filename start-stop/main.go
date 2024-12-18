package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func stopInstances(ctx context.Context) (string, error) {
	// Create a new session with default config
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)

	// Define filters
	filters := []*ec2.Filter{
		{
			Name:   aws.String("tag:Schedule"),
			Values: []*string{aws.String("WeekendStop")},
		},
		{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String("running")},
		},
	}

	// Describe instances with filters
	describeInput := &ec2.DescribeInstancesInput{Filters: filters}
	resp, err := svc.DescribeInstances(describeInput)
	if err != nil {
		return "", fmt.Errorf("error describing instances: %v", err)
	}

	// Extract instance IDs
	var instanceIDs []*string
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, instance.InstanceId)
		}
	}

	if len(instanceIDs) > 0 {
		// Stop instances
		stopInput := &ec2.StopInstancesInput{InstanceIds: instanceIDs}
		_, err := svc.StopInstances(stopInput)
		if err != nil {
			return "", fmt.Errorf("error stopping instances: %v", err)
		}
		return fmt.Sprintf("Stopped instances: %v", instanceIDs), nil
	}

	return "No instances to stop.", nil
}

func main() {
	lambda.Start(stopInstances)
}
