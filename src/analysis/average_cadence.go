package analysis

import (
	"bike/models"
	"bike/rider"

	log "github.com/sirupsen/logrus"
)

func AverageCadence(rider *rider.RIDER, ride *models.RideData) {
	log.Info("Average Cadence")
	var grandTotal float64
	var counter uint64

	for _, sample := range ride.Ride.Samples {
		if sample.Cad > 0 {
			grandTotal += sample.Cad
			counter++
		}
	}
	if counter == 0 || grandTotal == 0 {
		ride.Analysis.AverageCadence = 0
	} else {
		ac := grandTotal / float64(counter)
		ride.Analysis.AverageCadence = ac
	}

}
