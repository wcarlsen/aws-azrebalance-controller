package controller

import (
	"github.com/wcarlsen/aws-azrebalance-controller/internal/aws"
)

func observe(c aws.Clients, label string, nodegroupName string) (aws.Nodegroup, error) {
	// Get nodegroup details
	ng := aws.Nodegroup{Name: nodegroupName}
	err := ng.Get(c)
	if err != nil {
		return aws.Nodegroup{}, err
	}

	// Parse nodegroup labels
	err = ng.ParseLabels(label)
	if err != nil {
		return ng, err
	}

	// Loop over ASGs and get their details
	for i := range ng.Asgs {
		err := ng.Asgs[i].Get(c)
		if err != nil {
			return aws.Nodegroup{}, err
		}
	}
	return ng, nil
}
