package analysis

import (
	"bike/models"
	"bike/rider"

	log "github.com/sirupsen/logrus"
)

func Torque(rider *rider.RIDER, ride *models.RideData) {
	log.Info("Torque Power")

	for idx, sample := range ride.Ride.Samples {
		if sample.Watts > 0 {
			tq := sample.Watts / sample.Cad
			ride.Ride.Samples[idx].Torque = tq
		}
	}

}
