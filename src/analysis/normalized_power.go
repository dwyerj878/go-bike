package analysis

import (
	"bike/models"
	"bike/rider"
	"math"

	log "github.com/sirupsen/logrus"
)

func NormalizedPower(rider *rider.RIDER, ride *models.RideData) {
	log.Info("Normalized Power")
	length := len(ride.Ride.Samples) - 29
	if length <= 0 {
		log.Warn("Not enough samples for NP calculation")
		return
	}
	var grandTotal float64
	var counter uint64

	for i := range length {
		var total float64
		for j := range 30 {
			total += ride.Ride.Samples[i+j].Watts
		}
		vl := math.Pow(total/30.0, 4)
		grandTotal += vl
		counter++
	}
	np := grandTotal / float64(counter)
	np = math.Pow(np, 0.25)
	ride.Analysis.NormalizedPower = np

	ftp := float64(rider.Attributes[0].FTP)
	iff := np / ftp
	tss := (float64(ride.Ride.RecIntSecs) * np * iff) / (ftp * 3600) * 100
	ride.Analysis.IFF = iff
	ride.Analysis.TSS = tss

}
