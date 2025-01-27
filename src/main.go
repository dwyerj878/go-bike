package main

import (
	"bike/analysis"
	"bike/models"
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"

	"os"
)

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("Hello")
	log.Info(os.Args[1])
	fileName := os.Args[1]
	riderFileName := os.Args[2]
	rider, err := ReadRiderData(riderFileName)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	ride, err := ReadRide(fileName)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		analysis.SimpleAnalysis(rider, ride)
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		analysis.Temperature(ride)
	}(&wg)

	wg.Wait()
	//fmt.Println(result.Ride.Tags)
}

func ReadRide(fileName string) (*models.RIDE_DATA, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//var ride map[string]interface{}
	var ride models.RIDE_DATA
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&ride)
	if err != nil {
		return nil, err
	}
	log.Traceln(ride)
	return &ride, nil
}

func ReadRiderData(fileName string) (*models.RIDER, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var rider models.RIDER
	parser := json.NewDecoder(file)
	if err = parser.Decode(&rider); err != nil {
		return nil, err
	}
	maxHr := rider.Attributes[0].MaxHR

	rider.Attributes[0].HRZones = []models.RIDER_ZONE{
		{Min: 0, Max: uint32(float32(maxHr) * 0.5)},
		{Min: uint32(float32(maxHr)*0.5) + 1, Max: uint32(float32(maxHr) * 0.6)},
		{Min: uint32(float32(maxHr)*0.6) + 1, Max: uint32(float32(maxHr) * 0.7)},
		{Min: uint32(float32(maxHr)*0.7) + 1, Max: uint32(float32(maxHr) * 0.8)},
		{Min: uint32(float32(maxHr)*0.8) + 1, Max: uint32(float32(maxHr) * 0.9)},
		{Min: uint32(float32(maxHr)*0.9) + 1, Max: uint32(float32(maxHr) * 2.0)},
	}
	//fmt.Print(rider.Attributes[0].HRZones)

	ftp := rider.Attributes[0].FTP
	rider.Attributes[0].PowerZones = []models.RIDER_ZONE{
		{Min: 0, Max: uint32(float32(ftp) * 0.2)},
		{Min: uint32(float32(ftp)*0.2) + 1, Max: uint32(float32(ftp) * 0.5)},
		{Min: uint32(float32(ftp)*0.5) + 1, Max: uint32(float32(ftp) * 0.7)},
		{Min: uint32(float32(ftp)*0.7) + 1, Max: uint32(float32(ftp) * 0.85)},
		{Min: uint32(float32(ftp)*0.85) + 1, Max: uint32(float32(ftp) * 1.0)},
		{Min: uint32(float32(ftp)*1.0) + 1, Max: uint32(float32(ftp) * 1.15)},
		{Min: uint32(float32(ftp)*1.15) + 1, Max: uint32(float32(ftp) * 20.0)},
	}
	log.Traceln(rider)
	return &rider, nil
}
