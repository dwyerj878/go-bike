package analysis

import (
	"bike/models"
	"bike/rider"
	"math"

	log "github.com/sirupsen/logrus"
)

func ZoneIntervals(rider *rider.RIDER, ride *models.RideData) {
	log.Info("Intervals")
	var intervals []models.RideAnalysisZoneInterval

	var currentZone int
	var lastZone int
	var count uint32
	for _, sample := range ride.Ride.Samples {
		for idx, zone := range rider.Attributes[0].PowerZones {
			watts := uint32(math.Round(sample.Watts))
			if watts >= zone.Min && watts <= zone.Max {
				if idx != currentZone {
					if count < 5 {
						intervals[len(intervals)-1].Seconds += count
					} else if len(intervals) > 0 && lastZone == currentZone {
						intervals[len(intervals)-1].Seconds += count
					} else {
						intervals = append(intervals,
							models.RideAnalysisZoneInterval{
								Zone:    uint32(currentZone) + 1,
								Seconds: count,
							})
						lastZone = currentZone
					}
					currentZone = idx
					count = 1
				} else {
					count++
				}
				continue
			}
		}

	}
	if count > 0 {
		intervals = append(intervals,
			models.RideAnalysisZoneInterval{
				Zone:    uint32(currentZone) + 1,
				Seconds: uint32(count),
			})
	}
	ride.Analysis.ZoneIntervals = intervals
}
