package analysis

import (
	"bike/models"
	"bike/rider"

	log "github.com/sirupsen/logrus"
)

func FTPTimes(rider *rider.RIDER, ride *models.RideData) {
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

func PowerZoneTimes(rider *rider.RIDER, ride *models.RideData) {
	for zoneIdx, zone := range rider.Attributes[0].PowerZones {
		zoneData := models.RideAnalysisZone{Zone: uint8(zoneIdx + 1), Count: 0, Min: zone.Min, Max: zone.Max}
		ride.Analysis.PowerZones = append(ride.Analysis.PowerZones, zoneData)
	}
	var sampleCount uint64
	for _, sample := range ride.Ride.Samples {
		for zoneIdx, zone := range rider.Attributes[0].PowerZones {
			if sample.Watts >= float64(zone.Min) && sample.Watts <= float64(zone.Max) {
				ride.Analysis.PowerZones[zoneIdx].Count++
			}
		}
		sampleCount++
	}
	for idx, zone := range ride.Analysis.PowerZones {
		pct := float64(zone.Count) / float64(sampleCount) * 100
		ride.Analysis.PowerZones[idx].Percent = pct
	}
}

func Temperature(rider *rider.RIDER, ride *models.RideData) {
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

func MaxPower(rider *rider.RIDER, ride *models.RideData) {
	max := float64(0)
	for _, sample := range ride.Ride.Samples {
		if sample.Watts > max {
			max = sample.Watts
		}
	}
	log.Debugf("Max Watts :  Max : %f", max)
	ride.Analysis.MaxWatts = max
}

func HRZoneTimes(rider *rider.RIDER, ride *models.RideData) {
	for zoneIdx, zone := range rider.Attributes[0].HRZones {
		zoneData := models.RideAnalysisZone{Zone: uint8(zoneIdx + 1), Count: 0, Min: zone.Min, Max: zone.Max}
		ride.Analysis.HRZones = append(ride.Analysis.HRZones, zoneData)
	}
	var sampleCount uint64
	for _, sample := range ride.Ride.Samples {
		for zoneIdx, zone := range rider.Attributes[0].HRZones {
			if sample.Hr >= float64(zone.Min) && sample.Hr <= float64(zone.Max) {
				ride.Analysis.HRZones[zoneIdx].Count++
			}
		}
		sampleCount++
	}
	for idx, zone := range ride.Analysis.HRZones {
		pct := float64(zone.Count) / float64(sampleCount) * 100
		ride.Analysis.HRZones[idx].Percent = pct
	}
}
