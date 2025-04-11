package analysis

import (
	"bike/models"
	"bike/rider"
)

func ExecuteAnalysis(activeRider *rider.RIDER, ride *models.RIDE_DATA) {
	analysisFunctions := []func(*rider.RIDER, *models.RIDE_DATA){
		PowerZoneTimes,
		FTPTimes,
		Temperature,
		MaxPower,
		NormalizedPower,
		HRZoneTimes,
		AveragePower,
		AverageSpeed,
		// ZoneIntervals,
		Torque,
		AverageCadence,
	}
	for _, fnc := range analysisFunctions {
		fnc(activeRider, ride)
	}
}
