package controller

import (
	"github.com/wcarlsen/aws-azrebalance-controller/internal/aws"
)

func diff(ng aws.Nodegroup) []actType {
	var ds []actType
	for _, asg := range ng.Asgs {
		d := actType{asgName: asg.Name}
		for _, sp := range asg.SuspendedProcess {
			if *sp.ProcessName == aws.AZRebalance {
				if !ng.LabelBool {
					d.resume = true
				}
			} else {
				if ng.LabelBool {
					d.suspend = true
				}
			}
		}
		ds = append(ds, d)
	}
	return ds
}
