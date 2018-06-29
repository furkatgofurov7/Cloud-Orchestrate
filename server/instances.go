package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2 Response Json Object
type EC2Response struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	Monitoring string `json:"monitoring"`
	KeyName    string `json:"keyname"`
	// SecurityGroup    string `json:"sec_group"`
	// PrivateDNS       string `json:"private_dns"`
	// PrivateIP        string `json:"private_ip"`
	AvailabilityZone string `json:"availability_zone"`
}

func listInstances() []EC2Response {
	// Load session from shared config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new EC2 client
	ec2Svc := ec2.New(sess)

	result, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	response := []EC2Response{}
	for idx, res := range result.Reservations {
		fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range result.Reservations[idx].Instances {
			slot := EC2Response{*inst.InstanceId, *inst.InstanceType, *inst.State.Name,
				*inst.Monitoring.State, *inst.KeyName, //*inst.SecurityGroups[0].GroupName,
				//*inst.PrivateDnsName, *inst.PrivateIpAddress,
				*inst.Placement.AvailabilityZone}
			response = append(response, slot)
		}
	}
	return response

}

// Function to start and stop instances
func commandInstance(command string, instanceId string) string {
	// Load session from shared config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new EC2 client
	svc := ec2.New(sess)

	if command == "start" {
		input := &ec2.StartInstancesInput{
			InstanceIds: []*string{
				aws.String(instanceId),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.StartInstances(input)
		awsErr, ok := err.(awserr.Error)

		if ok && awsErr.Code() == "DryRunOperation" {
			// Let's now set dry run to be false. This will allow us to start the instances
			input.DryRun = aws.Bool(false)
			result, err = svc.StartInstances(input)
			if err != nil {
				fmt.Println("Error", err)
				return "Error"
			} else {
				fmt.Println("Success", result.StartingInstances)
				return "Started"
			}
		} else { // This could be due to a lack of permissions
			fmt.Println("Error", err)
			return "Error"
		}
	} else { // Turn instances off
		input := &ec2.StopInstancesInput{
			InstanceIds: []*string{
				aws.String(instanceId),
			},
			DryRun: aws.Bool(true),
		}
		result, err := svc.StopInstances(input)
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			result, err = svc.StopInstances(input)
			if err != nil {
				fmt.Println("Error", err)
				return "Error"
			} else {
				fmt.Println("Success", result.StoppingInstances)
				return "Stopped"
			}
		} else {
			fmt.Println("Error", err)
			return "Error"
		}
	}
}
