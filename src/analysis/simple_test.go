package analysis_test

import (
	"bike/analysis"
	"bike/models"
	"bike/rider"
	"math"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestTemperature(t *testing.T) {
	log.Debug("Test Temperature")
	ride := models.RideData{
		Ride: models.Ride{
			Samples: []models.RideSample{
				{Temp: 1},
				{Temp: 100},
				{Temp: 22},
			},
		},
	}
	analysis.Temperature(nil, &ride)
	if ride.Analysis.MaxTemp != 100 {
		t.Error("Max Temp")
	}
	if ride.Analysis.MinTemp != 1 {
		t.Error("Min Temp")
	}
}

func TestFTPTimes(t *testing.T) {
	log.Debug("Test FTP Times")
	ride := models.RideData{
		Ride: models.Ride{
			Samples: []models.RideSample{
				{Watts: 1},
				{Watts: 100},
				{Watts: 99},
				{Watts: 101},
				{Watts: 0},
				{Watts: 0.3},
			},
		},
	}
	rider := rider.RIDER{
		Attributes: []rider.RIDER_ATTRIBUTES{
			{
				FTP: 100,
			},
		},
	}
	analysis.FTPTimes(&rider, &ride)
	if ride.Analysis.FTP.Zero != 2 {
		t.Error("Incorrect Zero Count")
	}
	if ride.Analysis.FTP.Under != 2 {
		t.Error("Incorrect Under Count")
	}
	if ride.Analysis.FTP.Over != 2 {
		t.Error("Incorrect Over Count")
	}
}

func TestPowerZoneTimes(t *testing.T) {
	log.Debug("Test Zone Times")
	ride := models.RideData{
		Ride: models.Ride{
			Samples: []models.RideSample{
				{Watts: 1},
				{Watts: 11},
				{Watts: 30},
				{Watts: 33},
				{Watts: 44},
				{Watts: 55},
				{Watts: 100},
			},
		},
	}
	rider := rider.RIDER{
		Attributes: []rider.RIDER_ATTRIBUTES{
			{
				PowerZones: []rider.RIDER_ZONE{
					{Min: 0, Max: 10},
					{Min: 11, Max: 20},
					{Min: 21, Max: 30},
					{Min: 31, Max: 40},
					{Min: 41, Max: 50},
					{Min: 51, Max: 60},
					{Min: 61, Max: 100},
				},
			},
		},
	}
	analysis.PowerZoneTimes(&rider, &ride)
	for idx := 0; idx < 7; idx++ {
		if ride.Analysis.PowerZones[idx].Count != 1 {
			t.Errorf("Incorrect count for zone %d", idx+1)
		}
		expected := math.Trunc(100.0 / 7.0)
		if math.Trunc(float64(ride.Analysis.PowerZones[idx].Percent)) != expected {
			t.Errorf("Incorrect pct for zone %d %f expected %f", idx, ride.Analysis.PowerZones[idx].Percent, expected)
		}
	}

}

func TestMaxPower(t *testing.T) {
	log.Debug("Test FTP Times")
	ride := models.RideData{
		Ride: models.Ride{
			Samples: []models.RideSample{
				{Watts: 1},
				{Watts: 100},
				{Watts: 99},
				{Watts: 101},
				{Watts: 0},
				{Watts: 0.3},
			},
		},
	}
	analysis.MaxPower(nil, &ride)
	if ride.Analysis.MaxWatts != 101 {
		t.Error("Incorrect Max Watts")
	}
}

func TestHeartRateZoneTimes(t *testing.T) {
	log.Debug("Test Zone Times")
	ride := models.RideData{
		Ride: models.Ride{
			Samples: []models.RideSample{
				{Hr: 1},
				{Hr: 11},
				{Hr: 30},
				{Hr: 33},
				{Hr: 44},
				{Hr: 55},
				{Hr: 100},
			},
		},
	}
	rider := rider.RIDER{
		Attributes: []rider.RIDER_ATTRIBUTES{
			{
				HRZones: []rider.RIDER_ZONE{
					{Min: 0, Max: 10},
					{Min: 11, Max: 20},
					{Min: 21, Max: 30},
					{Min: 31, Max: 40},
					{Min: 41, Max: 50},
					{Min: 51, Max: 60},
					{Min: 61, Max: 100},
				},
			},
		},
	}
	analysis.HRZoneTimes(&rider, &ride)
	for idx := 0; idx < 7; idx++ {
		if ride.Analysis.HRZones[idx].Count != 1 {
			t.Errorf("Incorrect count for zone %d", idx+1)
		}
		expected := math.Trunc(100.0 / 7.0)
		if math.Trunc(float64(ride.Analysis.HRZones[idx].Percent)) != expected {
			t.Errorf("Incorrect pct for zone %d %f expected %f", idx, ride.Analysis.PowerZones[idx].Percent, expected)
		}
	}

}
