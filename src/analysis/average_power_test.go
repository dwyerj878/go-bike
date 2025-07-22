package analysis_test

import (
	"bike/analysis"
	"bike/models"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestAveragePower(t *testing.T) {
	log.Debug("Test Average Power")

	testCases := []struct {
		name          string
		samples       []models.RideSample
		expectedPower float64
	}{
		{
			name: "Valid Samples",
			samples: []models.RideSample{
				{Watts: 100},
				{Watts: 200},
				{Watts: 300},
				{Watts: 0},
				{Watts: 50},
			},
			expectedPower: 162.5,
		},
		{
			name:          "No Valid Samples",
			samples:       []models.RideSample{{Watts: 0}, {Watts: 0}, {Watts: 0}},
			expectedPower: 0,
		},
		{
			name: "Single Valid Sample",
			samples: []models.RideSample{
				{Watts: 0},
				{Watts: 150},
				{Watts: 0},
			},
			expectedPower: 150,
		},
		{
			name: "All valid",
			samples: []models.RideSample{
				{Watts: 100},
				{Watts: 200},
				{Watts: 300},
				{Watts: 50},
			},
			expectedPower: 162.5,
		},
		{
			name:          "Empty samples",
			samples:       []models.RideSample{},
			expectedPower: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ride := models.RideData{
				Ride: models.Ride{
					Samples: tc.samples,
				},
			}
			analysis.AveragePower(nil, &ride)
			if ride.Analysis.AveragePower != tc.expectedPower {
				t.Errorf("Expected average power to be %f, but got %f", tc.expectedPower, ride.Analysis.AveragePower)
			}
		})
	}
}
