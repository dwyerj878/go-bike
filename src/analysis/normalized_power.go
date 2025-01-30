package analysis

import (
	"bike/models"
	"bike/rider"
	"math"

	log "github.com/sirupsen/logrus"
)

func NormalizedPower(rider *rider.RIDER, ride *models.RIDE_DATA) {
	log.Info("Normalized Power")
	length := len(ride.Ride.Samples) - 29
	if length <= 0 {
		log.Warn("Not enough samples for NP calculation")
		return
	}
	var grandTotal float64
	var counter uint64

	for i := 0; i < length; i++ {
		var total float64
		for j := 0; j < 30; j++ {
			total += ride.Ride.Samples[i+j].Watts
		}
		vl := math.Pow(total/30.0, 4)
		grandTotal += vl
		counter++
	}
	np := grandTotal / float64(counter)
	np = math.Pow(np, 0.25)
	ride.Analysis.NP = np

}
