package controller

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	log "github.com/sirupsen/logrus"
	"github.com/wcarlsen/aws-azrebalance-controller/internal/aws"
)

const (
	action  = "AZRebalancing process %s"
	nothing = "Nothing to do AZRebalancing already %s"
)

type actType struct {
	asgName string
	resume  bool
	suspend bool
}

func resumeAZReblancing(c aws.Clients, asgName string) error {
	if c.DryRun {
		return nil
	}
	_, err := c.Asg.ResumeProcesses(c.Ctx, &autoscaling.ResumeProcessesInput{
		AutoScalingGroupName: &asgName,
		ScalingProcesses:     []string{aws.AZRebalance},
	})
	return err
}

func suspendAZRebalancing(c aws.Clients, asgName string) error {
	if c.DryRun {
		return nil
	}
	_, err := c.Asg.SuspendProcesses(c.Ctx, &autoscaling.SuspendProcessesInput{
		AutoScalingGroupName: &asgName,
		ScalingProcesses:     []string{aws.AZRebalance},
	})
	return err
}

func act(c aws.Clients, a actType, ng aws.Nodegroup) error {
	if a.resume {
		err := resumeAZReblancing(c, a.asgName)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"nodegroup": ng.Name,
			"asg":       a.asgName,
		}).Info("AZRebalancing process resumed")
		return nil
	}
	if a.suspend {
		err := suspendAZRebalancing(c, a.asgName)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"nodegroup": ng.Name,
			"asg":       a.asgName,
		}).Info("AZRebalancing process suspended")
		return nil
	}

	current := "resumed"
	if ng.LabelBool {
		current = "suspended"
	}
	log.WithFields(log.Fields{
		"nodegroup": ng.Name,
		"asg":       a.asgName,
	}).Info(fmt.Sprintf("Nothing to do AZRebalancing process already %s", current))
	return nil
}
