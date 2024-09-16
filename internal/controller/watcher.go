package controller

import (
	"github.com/aws/aws-sdk-go-v2/service/eks"
	log "github.com/sirupsen/logrus"
	"github.com/wcarlsen/aws-azrebalance-controller/internal/aws"
)

func watcher(c aws.Clients) (*eks.ListNodegroupsOutput, error) {
	// List nodegroups for cluster
	ngs, err := c.Eks.ListNodegroups(c.Ctx, &eks.ListNodegroupsInput{ClusterName: c.ClusterName})
	if err != nil {
		log.Error(err.Error())
		return &eks.ListNodegroupsOutput{}, err
	}
	return ngs, nil
}
