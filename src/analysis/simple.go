package analysis

import (
	"bike/models"

	log "github.com/sirupsen/logrus"
)

func SimpleAnalysis(rider *models.RIDER, ride *models.RIDE_DATA) {
	var max float32
	var powerRanges [100]uint64
	ftp := float32(rider.Attributes[0].FTP)
	wattRange := float32(rider.Attributes[0].FTP / 10)
	overCount := uint64(0)
	underCount := uint64(0)
	zeroCount := uint64(0)
	for _, sample := range ride.Ride.Samples {
		if sample.Watts > max {
			max = sample.Watts
		}
		if sample.Watts <= 1.0 {
			zeroCount++
		} else if sample.Watts >= ftp {
			overCount++
		} else {
			underCount++
		}
		idx := int16(sample.Watts / wattRange)
		powerRanges[idx] = powerRanges[idx] + 1
	}
	maxIdx := uint16(max / 25)
	var i uint16
	for i = 0; i < maxIdx; i++ {
		log.Debugf("range %d %d : count %d\n", i*25, i*25+25, powerRanges[i])
	}

	log.Debugf("max : %f\n", max)
	log.Debugf("Zero %d Over %d Under %d\n", zeroCount, overCount, underCount)
}

func Temperature(rider *models.RIDER, ride *models.RIDE_DATA) {
	min := float32(500)
	max := float32(0)
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
