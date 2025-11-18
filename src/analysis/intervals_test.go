package analysis

import (
	"bike/models"
	"bike/rider"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPowerZoneIndex(t *testing.T) {
	zones := []rider.Zone{
		{Min: 0, Max: 100},
		{Min: 101, Max: 150},
		{Min: 151, Max: 200},
	}

	testCases := []struct {
		name          string
		watts         uint32
		expectedIndex int
	}{
		{"In Zone 1", 50, 0},
		{"Lower bound Zone 2", 101, 1},
		{"Upper bound Zone 2", 150, 1},
		{"In Zone 3", 175, 2},
		{"Below all zones", 0, 0}, // Assuming 0 is part of the first zone
		{"Above all zones", 250, -1},
		{"Between zones", 100, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedIndex, getPowerZoneIndex(tc.watts, zones))
		})
	}
}

func TestZoneIntervals(t *testing.T) {
	testRider := &rider.Rider{
		Attributes: []rider.Attributes{
			{
				PowerZones: []rider.Zone{
					{Min: 0, Max: 100},   // Zone 1 (index 0)
					{Min: 101, Max: 150}, // Zone 2 (index 1)
					{Min: 151, Max: 200}, // Zone 3 (index 2)
				},
			},
		},
	}

	testCases := []struct {
		name              string
		ride              *models.RideData
		rider             *rider.Rider
		expectedIntervals []models.RideAnalysisZoneInterval
	}{
		{
			name:              "No ride data",
			ride:              nil,
			rider:             testRider,
			expectedIntervals: nil,
		},
		{
			name:              "No rider data",
			ride:              &models.RideData{},
			rider:             nil,
			expectedIntervals: nil,
		},
		{
			name:              "Rider with no attributes",
			ride:              &models.RideData{Ride: models.Ride{Samples: []models.RideSample{{Watts: 120}}}},
			rider:             &rider.Rider{},
			expectedIntervals: nil,
		},
		{
			name:              "Ride with no samples",
			ride:              &models.RideData{Ride: models.Ride{Samples: []models.RideSample{}}},
			rider:             testRider,
			expectedIntervals: []models.RideAnalysisZoneInterval{},
		},
		{
			name: "Single valid interval",
			ride: &models.RideData{
				Ride: models.Ride{
					RecIntSecs: 1,
					Samples: []models.RideSample{
						{Watts: 50}, {Watts: 50},
						{Watts: 120}, {Watts: 125}, {Watts: 130}, {Watts: 128}, {Watts: 135}, // 5 samples in zone 2
						{Watts: 50}, {Watts: 50},
					},
				},
			},
			rider: testRider,
			expectedIntervals: []models.RideAnalysisZoneInterval{
				{Zone: 2, Seconds: 5},
			},
		},
		{
			name: "Interval too short",
			ride: &models.RideData{
				Ride: models.Ride{
					RecIntSecs: 1,
					Samples: []models.RideSample{
						{Watts: 50}, {Watts: 50},
						{Watts: 120}, {Watts: 125}, {Watts: 130}, {Watts: 128}, // 4 samples
						{Watts: 50}, {Watts: 50},
					},
				},
			},
			rider:             testRider,
			expectedIntervals: []models.RideAnalysisZoneInterval{},
		},
		{
			name: "Multiple valid intervals",
			ride: &models.RideData{
				Ride: models.Ride{
					RecIntSecs: 1,
					Samples: []models.RideSample{
						{Watts: 110}, {Watts: 110}, {Watts: 110}, {Watts: 110}, {Watts: 110}, // 5 in zone 2
						{Watts: 50}, {Watts: 50},
						{Watts: 160}, {Watts: 160}, {Watts: 160}, {Watts: 160}, {Watts: 160}, {Watts: 160}, // 6 in zone 3
						{Watts: 50},
					},
				},
			},
			rider: testRider,
			expectedIntervals: []models.RideAnalysisZoneInterval{
				{Zone: 2, Seconds: 5},
				{Zone: 3, Seconds: 6},
			},
		},
		{
			name: "Interval at the end of the ride",
			ride: &models.RideData{
				Ride: models.Ride{
					RecIntSecs: 1,
					Samples: []models.RideSample{
						{Watts: 50}, {Watts: 50}, {Watts: 50},
						{Watts: 180}, {Watts: 180}, {Watts: 180}, {Watts: 180}, {Watts: 180}, {Watts: 180}, // 6 in zone 3
					},
				},
			},
			rider: testRider,
			expectedIntervals: []models.RideAnalysisZoneInterval{
				{Zone: 3, Seconds: 6},
			},
		},
		{
			name: "No consecutive samples in a zone",
			ride: &models.RideData{
				Ride: models.Ride{
					RecIntSecs: 1,
					Samples: []models.RideSample{
						{Watts: 50}, {Watts: 120}, {Watts: 50}, {Watts: 120}, {Watts: 50}, {Watts: 120}, {Watts: 50},
					},
				},
			},
			rider:             testRider,
			expectedIntervals: []models.RideAnalysisZoneInterval{},
		},
		{
			name: "RecIntSecs is not 1",
			ride: &models.RideData{
				Ride: models.Ride{
					RecIntSecs: 2,
					Samples: []models.RideSample{
						{Watts: 120}, {Watts: 125}, {Watts: 130}, {Watts: 128}, {Watts: 135}, // 5 samples in zone 2
						{Watts: 50},
					},
				},
			},
			rider: testRider,
			expectedIntervals: []models.RideAnalysisZoneInterval{
				{Zone: 2, Seconds: 10}, // 5 samples * 2 seconds
			},
		},
		{
			name: "RecIntSecs is 0 (defaults to 1)",
			ride: &models.RideData{
				Ride: models.Ride{
					RecIntSecs: 0,
					Samples: []models.RideSample{
						{Watts: 120}, {Watts: 125}, {Watts: 130}, {Watts: 128}, {Watts: 135}, // 5 samples in zone 2
						{Watts: 50},
					},
				},
			},
			rider: testRider,
			expectedIntervals: []models.RideAnalysisZoneInterval{
				{Zone: 2, Seconds: 5}, // 5 samples * 1 second
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// The function modifies the ride in place, so we don't need to capture a return value.
			ZoneIntervals(tc.rider, tc.ride)

			if tc.ride == nil {
				// If ride is nil, no further checks can be done. The test is to ensure no panic.
				return
			}

			// For cases where intervals are expected to be nil or empty
			if tc.expectedIntervals == nil {
				assert.Nil(t, tc.ride.Analysis.ZoneIntervals)
			} else {
				assert.ElementsMatch(t, tc.expectedIntervals, tc.ride.Analysis.ZoneIntervals)
			}
		})
	}
}
