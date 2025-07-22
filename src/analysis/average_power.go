package analysis

import (
	"bike/models"
	"bike/rider"

	log "github.com/sirupsen/logrus"
)

func AveragePower(rider *rider.RIDER, ride *models.RideData) {
	log.Info("Average Power")
	var grandTotal float64
	var counter uint64

	for _, sample := range ride.Ride.Samples {
		if sample.Watts > 0 {
			grandTotal += sample.Watts
			counter++
		}
	}
	if counter == 0 || grandTotal == 0 {
		ride.Analysis.AveragePower = 0
	} else {
		ap := grandTotal / float64(counter)
		ride.Analysis.AveragePower = ap
	}

}
