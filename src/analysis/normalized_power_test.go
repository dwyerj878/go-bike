package analysis

import (
	"bike/models"
	"bike/rider"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedPower(t *testing.T) {
	// Test case 1: Constant power
	t.Run("Constant Power", func(t *testing.T) {
		rider := &rider.Rider{
			Attributes: []rider.Attributes{
				{FTP: 250},
			},
		}
		ride := &models.RideData{
			Ride: models.Ride{
				Samples:    make([]models.RideSample, 60), // 60 samples
				RecIntSecs: 60,
			},
			Analysis: models.RideAnalysis{},
		}

		for i := range ride.Ride.Samples {
			ride.Ride.Samples[i].Watts = 200
		}

		NormalizedPower(rider, ride)

		assert.InDelta(t, 200, ride.Analysis.NormalizedPower, 0.01, "Normalized Power should be 200 for constant power")

		// IFF = NP / FTP = 200 / 250 = 0.8
		assert.InDelta(t, 0.8, ride.Analysis.IFF, 0.01, "IFF calculation is incorrect")

		// TSS = (seconds * NP * IF) / (FTP * 3600) * 100
		expectedTSS := (float64(60) * 200 * 0.8) / (250 * 3600) * 100
		assert.InDelta(t, expectedTSS, ride.Analysis.TSS, 0.01, "TSS calculation is incorrect")
	})

	// Test case 2: Not enough samples
	t.Run("Not Enough Samples", func(t *testing.T) {
		rider := &rider.Rider{
			Attributes: []rider.Attributes{
				{FTP: 250},
			},
		}
		ride := &models.RideData{
			Ride: models.Ride{
				Samples: make([]models.RideSample, 20), // 20 samples, less than 30
			},
			Analysis: models.RideAnalysis{},
		}

		for i := range ride.Ride.Samples {
			ride.Ride.Samples[i].Watts = 200
		}

		NormalizedPower(rider, ride)

		// NP should be 0 as it's not calculated
		assert.Equal(t, float64(0), ride.Analysis.NormalizedPower, "Normalized Power should be 0 for insufficient samples")
		assert.Equal(t, float64(0), ride.Analysis.IFF, "IFF should be 0 for insufficient samples")
		assert.Equal(t, float64(0), ride.Analysis.TSS, "TSS should be 0 for insufficient samples")
	})

	// Test case 3: Variable power
	t.Run("Variable Power", func(t *testing.T) {
		rider := &rider.Rider{
			Attributes: []rider.Attributes{
				{FTP: 200},
			},
		}
		ride := &models.RideData{
			Ride: models.Ride{
				Samples:    make([]models.RideSample, 31), // 31 samples
				RecIntSecs: 31,
			},
			Analysis: models.RideAnalysis{},
		}

		// Alternate between 150W and 250W. The average is 200W.
		// The 30-second rolling average will always be ~200W.
		for i := range ride.Ride.Samples {
			if i%2 == 0 {
				ride.Ride.Samples[i].Watts = 150
			} else {
				ride.Ride.Samples[i].Watts = 250
			}
		}

		NormalizedPower(rider, ride)

		// The NP will be slightly different from the average due to the 4th power calculation,
		// but for this stable alternating pattern, it will be very close to the average.
		// Let's calculate the expected NP more accurately.
		// Window 1 (0-29): 15*150W, 15*250W. Avg = 200.
		// Window 2 (1-30): 15*250W, 15*150W. Avg = 200.
		// So, the average of the 4th powers of the rolling averages is 200^4.
		// The 4th root of that is 200.
		assert.InDelta(t, 200, ride.Analysis.NormalizedPower, 0.01, "Normalized Power calculation is incorrect for variable power")

		// IFF = NP / FTP = 200 / 200 = 1.0
		assert.InDelta(t, 1.0, ride.Analysis.IFF, 0.01, "IFF calculation is incorrect")

		// TSS = (seconds * NP * IF) / (FTP * 3600) * 100
		expectedTSS := (float64(31) * 200 * 1.0) / (200 * 3600) * 100
		assert.InDelta(t, expectedTSS, ride.Analysis.TSS, 0.01, "TSS calculation is incorrect")
	})
}
