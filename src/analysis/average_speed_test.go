package analysis_test

import (
	"bike/analysis"
	"bike/models"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestAverageSpeed(t *testing.T) {
	log.Debug("Test Average Speed")

	testCases := []struct {
		name          string
		samples       []models.RIDE_SAMPLE
		expectedSpeed float64
	}{
		{
			name: "Valid Samples",
			samples: []models.RIDE_SAMPLE{
				{Kph: 10},
				{Kph: 20},
				{Kph: 30},
				{Kph: 0},
				{Kph: 5},
			},
			expectedSpeed: 16.25,
		},
		{
			name:          "No Valid Samples",
			samples:       []models.RIDE_SAMPLE{{Kph: 0}, {Kph: 0}, {Kph: 0}},
			expectedSpeed: 0,
		},
		{
			name: "Single Valid Sample",
			samples: []models.RIDE_SAMPLE{
				{Kph: 0},
				{Kph: 15},
				{Kph: 0},
			},
			expectedSpeed: 15,
		},
		{
			name: "All valid",
			samples: []models.RIDE_SAMPLE{
				{Kph: 10},
				{Kph: 20},
				{Kph: 30},
				{Kph: 5},
			},
			expectedSpeed: 16.25,
		},
		{
			name:          "Empty samples",
			samples:       []models.RIDE_SAMPLE{},
			expectedSpeed: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ride := models.RIDE_DATA{
				Ride: models.RIDE{
					Samples: tc.samples,
				},
			}
			analysis.AverageSpeed(nil, &ride)
			if ride.Analysis.AverageSpeed != tc.expectedSpeed {
				t.Errorf("Expected average speed to be %f, but got %f", tc.expectedSpeed, ride.Analysis.AverageSpeed)
			}
		})
	}
}
