package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

const AZRebalance string = "AZRebalance"

type Asg struct {
	Name             string
	SuspendedProcess []types.SuspendedProcess
}

func (a *Asg) Get(c Clients) error {
	// Describe asg
	dasg, err := c.Asg.DescribeAutoScalingGroups(c.Ctx, &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{a.Name},
	})
	if err != nil {
		return err
	}
	for _, asg := range dasg.AutoScalingGroups {
		a.SuspendedProcess = asg.SuspendedProcesses
	}
	return nil
}
