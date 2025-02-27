package analysis_test

import (
	"bike/analysis"
	"bike/models"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestTorque(t *testing.T) {
	log.Debug("Test Torque")
	ride := models.RIDE_DATA{
		Ride: models.RIDE{
			Samples: []models.RIDE_SAMPLE{
				{Watts: 270, Cad: 90},
				{Watts: 100, Cad: 100},
				{Watts: 600, Cad: 120},
			},
		},
	}
	analysis.Torque(nil, &ride)
	if ride.Ride.Samples[0].Torque != 3.0 {
		t.Errorf("Incorrect torque calculation %f", ride.Ride.Samples[0].Torque)
	}

}
