package controller

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wcarlsen/aws-azrebalance-controller/internal/aws"
)

func diff(ng aws.Nodegroup) []actType {
	var ds []actType
	for _, asg := range ng.Asgs {
		d := actType{asgName: asg.Name}
		for _, sp := range asg.SuspendedProcess {
			if *sp.ProcessName == aws.AZRebalance {
				if !ng.LabelBool {
					if asg.Instances > 0 {
						log.WithFields(log.Fields{
							"nodegroup": ng.Name,
							"asg":       asg.Name,
						}).Info(fmt.Sprintf("AZRebalance process should be resumed, but instances running is %d so it will be suspended", asg.Instances))
						d.suspend = true
					} else {
						d.resume = true
					}
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
