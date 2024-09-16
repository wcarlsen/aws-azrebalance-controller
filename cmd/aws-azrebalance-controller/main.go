package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	log "github.com/sirupsen/logrus"
	"github.com/wcarlsen/aws-azrebalance-controller/internal/aws"
	"github.com/wcarlsen/aws-azrebalance-controller/internal/controller"
)

const errGoRoutineStopped = "go routine for reconciling stopped"

var (
	region         = flag.String("region", os.Getenv("AWS_REGION"), "AWS region")
	clusterName    = flag.String("cluster-name", os.Getenv("CLUSTER_NAME"), "Cluster name")
	nodegroupLabel = flag.String("label", os.Getenv("LABEL"), "Nodegroup target label")
	period         = flag.Duration("period", 10*time.Second, "Reconciliation period in seconds")
	debug          = flag.Bool("debug", false, "Enable debug logging")
	dryRun         = flag.Bool("dry-run", false, "Changes disabled")
)

func init() {
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	lvl := log.InfoLevel
	if *debug {
		lvl = log.DebugLevel
	}
	log.SetLevel(lvl)
}

func main() {

	// Setup
	log.WithFields(log.Fields{
		"region":       region,
		"cluster-name": clusterName,
		"label":        nodegroupLabel,
		"period":       period.Seconds(),
		"dry-run":      dryRun,
	}).Info("Starting AZ rebalancing controller")

	// Setup clients
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region))
	if err != nil {
		log.Fatal(err.Error())
	}

	clients := aws.Clients{
		Ctx:         ctx,
		Eks:         eks.NewFromConfig(cfg),
		Asg:         autoscaling.NewFromConfig(cfg),
		ClusterName: clusterName,
		DryRun:      *dryRun,
	}

	// Reconcile on period
	cancel := make(chan struct{})
	ticker := time.NewTicker(*period)
	defer ticker.Stop()

	go func() {
		for {
			controller.Reconile(clients, *nodegroupLabel)
			select {
			case <-ticker.C:
				continue
			case <-cancel:
				log.Fatal(errGoRoutineStopped)
			}
		}
	}()

	// Blocks until a signal is received (e.g. Ctrl+C).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

}
