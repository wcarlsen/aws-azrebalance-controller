package controller

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wcarlsen/aws-azrebalance-controller/internal/aws"
)

func Reconile(c aws.Clients, label string, instanceAware bool) {
	st := time.Now()
	log.Info("Reconciliation loop started")

	// Watch
	nodegroups, err := watcher(c)
	if err != nil {
		log.Fatal("Cannot list nodegroups")
	}

	var wg sync.WaitGroup

	for _, nodegroupName := range nodegroups.Nodegroups {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Observe
			ng, err := observe(c, label, nodegroupName)
			if err != nil {
				if err.Error() == aws.ErrLabelNotPresent {
					log.WithFields(log.Fields{
						"nodegroup": ng.Name,
					}).Debug("Skipping nodegroup matching label not found")
					return
				}
				log.WithFields(log.Fields{
					"nodegroup": ng.Name,
				}).Error(err.Error())
				return
			}

			// Diff
			acts := diff(ng, instanceAware)

			// Act
			for _, a := range acts {
				err := act(c, a, ng)
				if err != nil {
					log.WithFields(log.Fields{
						"nodegroup": ng.Name,
						"asg":       a.asgName,
					}).Error(err.Error())
				}
			}
		}()
	}
	wg.Wait()
	log.WithFields(log.Fields{
		"duration": time.Since(st).Seconds(),
	}).Info("Reconciliation loop ended")
}
