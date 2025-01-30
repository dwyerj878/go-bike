package analysis

import (
	"bike/models"
	"bike/rider"

	log "github.com/sirupsen/logrus"
)

func FTPTimes(rider *rider.RIDER, ride *models.RIDE_DATA) {
	ftp := float64(rider.Attributes[0].FTP)
	over := uint64(0)
	under := uint64(0)
	zero := uint64(0)
	for _, sample := range ride.Ride.Samples {
		if sample.Watts < 1.0 {
			zero++
		} else if sample.Watts >= ftp {
			over++
		} else {
			under++
		}
	}
	ride.Analysis.FTP.FTP = rider.Attributes[0].FTP
	ride.Analysis.FTP.Over = over
	ride.Analysis.FTP.Under = under
	ride.Analysis.FTP.Zero = zero
}

func ZoneTimes(rider *rider.RIDER, ride *models.RIDE_DATA) {
	for zoneIdx, _ := range rider.Attributes[0].PowerZones {
		ride.Analysis.ZONES = append(ride.Analysis.ZONES, models.RIDE_ANALYSIS_ZONE{Zone: uint8(zoneIdx + 1), Count: 0})
	}
	var sampleCount uint64
	for _, sample := range ride.Ride.Samples {
		for zoneIdx, zone := range rider.Attributes[0].PowerZones {
			if sample.Watts >= float64(zone.Min) && sample.Watts <= float64(zone.Max) {
				ride.Analysis.ZONES[zoneIdx].Count++
			}
		}
		sampleCount++
	}
	for idx, zone := range ride.Analysis.ZONES {
		pct := float64(zone.Count) / float64(sampleCount) * 100
		ride.Analysis.ZONES[idx].Percent = pct
	}

}

func Temperature(rider *rider.RIDER, ride *models.RIDE_DATA) {
	min := float64(500)
	max := float64(0)
	for _, sample := range ride.Ride.Samples {
		if sample.Temp > max {
			max = sample.Temp
		}
		if sample.Temp < min {
			min = sample.Temp
		}
	}
	log.Debugf("Tempt Min : %f  Max : %f", min, max)
	ride.Analysis.MaxTemp = max
	ride.Analysis.MinTemp = min
}

func MaxPower(rider *rider.RIDER, ride *models.RIDE_DATA) {
	max := float64(0)
	for _, sample := range ride.Ride.Samples {
		if sample.Watts > max {
			max = sample.Watts
		}
	}
	log.Debugf("Max Watts :  Max : %f", max)
	ride.Analysis.MaxWatts = max
}
