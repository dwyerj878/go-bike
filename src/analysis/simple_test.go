package analysis

import (
	"bike/models"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestTemperature(t *testing.T) {
	log.Debug("Test Temperature")
	ride := models.RIDE_DATA{
		Ride: models.RIDE{
			Samples: []models.RIDE_SAMPLE{
				{Temp: 1},
				{Temp: 100},
				{Temp: 22},
			},
		},
	}
	Temperature(nil, &ride)
	if ride.Analysis.MaxTemp != 100 {
		t.Error("Max Temp")
	}
	if ride.Analysis.MinTemp != 1 {
		t.Error("Min Temp")
	}
}
