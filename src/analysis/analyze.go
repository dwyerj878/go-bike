package analysis

import (
	"bike/models"
	"bike/rider"
	"sync"
)

func ExecuteAnalysis(activeRider *rider.RIDER, ride *models.RIDE_DATA) {
	var wg sync.WaitGroup
	analysisFunctions := []func(*rider.RIDER, *models.RIDE_DATA){
		PowerZoneTimes,
		FTPTimes,
		Temperature,
		MaxPower,
		NormalizedPower,
		HRZoneTimes,
		AveragePower,
		AverageSpeed,
		ZoneIntervals,
		Torque,
		AverageCadence,
	}
	for _, fnc := range analysisFunctions {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			fnc(activeRider, ride)
		}(&wg)

	}
	wg.Wait()
}
