package analysis

import (
	"bike/models"
	"bike/rider"

	log "github.com/sirupsen/logrus"
)

func AverageSpeed(rider *rider.RIDER, ride *models.RIDE_DATA) {
	log.Info("Average Speed")
	var grandTotal float64
	var counter uint64

	for _, sample := range ride.Ride.Samples {
		if sample.Kph > 0 {
			grandTotal += sample.Kph
			counter++
		}
	}
	averageSpeed := grandTotal / float64(counter)
	ride.Analysis.AverageSpeed = averageSpeed

}
