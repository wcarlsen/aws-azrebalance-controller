package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

type Clients struct {
	Eks         *eks.Client
	Asg         *autoscaling.Client
	Ctx         context.Context
	ClusterName *string
	DryRun      bool
}
