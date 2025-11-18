package analysis

import (
	"bike/models"
	"bike/rider"
)

func ExecuteAnalysis(activeRider *rider.Rider, ride *models.RideData) {
	analysisFunctions := []func(*rider.Rider, *models.RideData){
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
