package analysis_test

import (
	"bike/analysis"
	"bike/models"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestAverageCadence(t *testing.T) {
	log.Debug("Test Average Cadence")

	testCases := []struct {
		name            string
		samples         []models.RIDE_SAMPLE
		expectedCadence float64
	}{
		{
			name: "Valid Samples",
			samples: []models.RIDE_SAMPLE{
				{Cad: 100},
				{Cad: 200},
				{Cad: 300},
				{Cad: 0},
				{Cad: 50},
			},
			expectedCadence: 162.5,
		},
		{
			name:            "No Valid Samples",
			samples:         []models.RIDE_SAMPLE{{Cad: 0}, {Cad: 0}, {Cad: 0}},
			expectedCadence: 0,
		},
		{
			name: "Single Valid Sample",
			samples: []models.RIDE_SAMPLE{
				{Cad: 0},
				{Cad: 150},
				{Cad: 0},
			},
			expectedCadence: 150,
		},
		{
			name: "All valid",
			samples: []models.RIDE_SAMPLE{
				{Cad: 100},
				{Cad: 200},
				{Cad: 300},
				{Cad: 50},
			},
			expectedCadence: 162.5,
		},
		{
			name:            "Empty samples",
			samples:         []models.RIDE_SAMPLE{},
			expectedCadence: 0,
		},
		{
			name: "Some negatives samples",
			samples: []models.RIDE_SAMPLE{
				{Cad: -100},
				{Cad: 200},
				{Cad: 300},
				{Cad: 50},
			},
			expectedCadence: 183.33333333333334,
		},
		{
			name: "Negative samples",
			samples: []models.RIDE_SAMPLE{
				{Cad: -100},
				{Cad: -200},
				{Cad: -300},
				{Cad: -50},
			},
			expectedCadence: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ride := models.RIDE_DATA{
				Ride: models.RIDE{
					Samples: tc.samples,
				},
			}
			analysis.AverageCadence(nil, &ride)
			if ride.Analysis.AverageCadence != tc.expectedCadence {
				t.Errorf("Expected average cadence to be %f, but got %f", tc.expectedCadence, ride.Analysis.AverageCadence)
			}
		})
	}
}
