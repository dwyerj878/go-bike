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
	activeRider, err := rider.ReadRiderData(riderFileName)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	ride, err := ReadRide(fileName)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	t := []func(*rider.RIDER, *models.RIDE_DATA){
		analysis.ZoneTimes,
		analysis.FTPTimes,
		analysis.Temperature,
		analysis.MaxPower,
	}
	for _, fnc := range t {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			fnc(activeRider, ride)
		}(&wg)

	}
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
