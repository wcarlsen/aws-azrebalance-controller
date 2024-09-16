package aws

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/pkg/errors"
)

const (
	ErrLabelNotPresent string = "Nodegroup label is not present"
)

type Nodegroup struct {
	Name      string
	labels    map[string]string
	LabelBool bool
	Asgs      []Asg
}

func (ng *Nodegroup) Get(c Clients) error {
	// Describe nodegroup
	dng, err := c.Eks.DescribeNodegroup(c.Ctx, &eks.DescribeNodegroupInput{
		ClusterName:   c.ClusterName,
		NodegroupName: &ng.Name,
	})
	if err != nil {
		return err
	}
	ng.labels = dng.Nodegroup.Labels

	for _, autoscalingGroup := range dng.Nodegroup.Resources.AutoScalingGroups {
		asg := Asg{Name: *autoscalingGroup.Name}
		ng.Asgs = append(ng.Asgs, asg)
	}
	return nil
}

func (ng *Nodegroup) ParseLabels(label string) error {
	l, ok := ng.labels[label]
	if !ok {
		return errors.New(ErrLabelNotPresent)
	}
	lb, err := strconv.ParseBool(l)
	if err != nil {
		return err
	}
	ng.LabelBool = lb
	return nil
}
