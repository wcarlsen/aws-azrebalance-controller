package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/pkg/errors"
)

const AZRebalance string = "AZRebalance"

const errGotMoreASGsThanRequested string = "Describe single ASG returned more than one ASG"

type Asg struct {
	Name             string
	SuspendedProcess []string
	Instances        int
}

func (a *Asg) Get(c Clients) error {
	// Describe asg
	dasg, err := c.Asg.DescribeAutoScalingGroups(c.Ctx, &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{a.Name},
	})
	if err != nil {
		return err
	}
	if len(dasg.AutoScalingGroups) > 1 {
		return errors.New(errGotMoreASGsThanRequested)
	}
	asg := dasg.AutoScalingGroups[0] // we know we request only one
	a.Instances = len(asg.Instances)

	for _, sp := range asg.SuspendedProcesses {
		a.SuspendedProcess = append(a.SuspendedProcess, *sp.ProcessName)
	}
	return nil
}
