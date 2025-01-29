package main

import (
	"bike/analysis"
	"bike/models"
	"bike/rider"
	"encoding/json"
	"fmt"
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
	rider, err := rider.ReadRiderData(riderFileName)
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
		analysis.ZoneTimes(rider, ride)
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		analysis.Temperature(rider, ride)
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		analysis.FTPTimes(rider, ride)
	}(&wg)

	wg.Wait()
	//fmt.Println(result.Ride.Tags)
	b, err := json.MarshalIndent(ride.Analysis, "", "  ")
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("Analysis : %s", b)
		fmt.Printf("Analysis : %s", b)
	}
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
