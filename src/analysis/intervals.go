package analysis

import (
	"bike/models"
	"bike/rider"
	"math"
)

// A struct to hold interval data during processing
type powerInterval struct {
	startSampleIndex int
	endSampleIndex   int
	zoneIndex        int
}

// ZoneIntervals detects consecutive periods of time spent in a specific power zone.
// An interval is defined as a minimum of 5 consecutive samples within a single power zone.
func ZoneIntervals(theRider *rider.Rider, ride *models.RideData) {
	if ride == nil || len(ride.Ride.Samples) == 0 {
		return
	}

	// Assuming the latest rider attributes are the ones to use.
	if theRider == nil || len(theRider.Attributes) == 0 {
		return
	}
	powerZones := theRider.Attributes[0].PowerZones
	samples := ride.Ride.Samples
	recIntSecs := ride.Ride.RecIntSecs
	if recIntSecs == 0 {
		recIntSecs = 1 // Assume 1 second per sample if not specified
	}

	var foundIntervals []powerInterval
	var currentInterval *powerInterval

	for i, sample := range samples {
		sampleWatts := uint32(math.Round(sample.Watts))
		zoneIdx := getPowerZoneIndex(sampleWatts, powerZones)

		if currentInterval == nil {
			if zoneIdx != -1 {
				// Start of a potential new interval
				currentInterval = &powerInterval{
					startSampleIndex: i,
					zoneIndex:        zoneIdx,
				}
			}
		} else {
			// We are in an interval. Check if it continues.
			if zoneIdx != currentInterval.zoneIndex {
				// Interval has ended at the previous sample.
				currentInterval.endSampleIndex = i - 1

				// Check if it meets the minimum length requirement.
				if (currentInterval.endSampleIndex - currentInterval.startSampleIndex + 1) >= 5 {
					foundIntervals = append(foundIntervals, *currentInterval)
				}

				// Reset and check if a new interval starts with the current sample.
				if zoneIdx != -1 {
					currentInterval = &powerInterval{
						startSampleIndex: i,
						zoneIndex:        zoneIdx,
					}
				} else {
					currentInterval = nil
				}
			}
		}
	}

	// Handle the case where an interval is ongoing at the end of the ride.
	if currentInterval != nil {
		currentInterval.endSampleIndex = len(samples) - 1
		if (currentInterval.endSampleIndex - currentInterval.startSampleIndex + 1) >= 5 {
			foundIntervals = append(foundIntervals, *currentInterval)
		}
	}

	// Convert the found intervals into the format expected by RideAnalysis.
	analysisIntervals := make([]models.RideAnalysisZoneInterval, 0, len(foundIntervals))
	for _, interval := range foundIntervals {
		duration := uint32((interval.endSampleIndex - interval.startSampleIndex + 1) * recIntSecs)
		analysisIntervals = append(analysisIntervals, models.RideAnalysisZoneInterval{
			Zone:    uint32(interval.zoneIndex + 1), // Zones are 1-based
			Seconds: duration,
		})
	}

	ride.Analysis.ZoneIntervals = analysisIntervals
}

// getPowerZoneIndex finds which power zone a given wattage falls into.
// Returns the 0-based index of the zone, or -1 if not in any zone.
func getPowerZoneIndex(watts uint32, zones []rider.Zone) int {
	for i, zone := range zones {
		if watts >= zone.Min && watts <= zone.Max {
			return i
		}
	}
	return -1
}
