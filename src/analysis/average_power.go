package analysis

import (
	"bike/models"
	"bike/rider"

	log "github.com/sirupsen/logrus"
)

func AveragePower(rider *rider.RIDER, ride *models.RIDE_DATA) {
	log.Info("Average Power")
	var grandTotal float64
	var counter uint64

	for _, sample := range ride.Ride.Samples {
		if sample.Watts > 0 {
			grandTotal += sample.Watts
			counter++
		}
	}
	ap := grandTotal / float64(counter)
	ride.Analysis.AveragePower = ap

}
