package main

import (
	"bike/analysis"
	"bike/models"
	"bike/rider"
	"encoding/json"
	"fmt"

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

	ride, err := models.Read(fileName)
	if err != nil {
		panic(err)
	}
	analysis.ExecuteAnalysis(activeRider, ride)
	b, err := json.MarshalIndent(ride.Analysis, "", "  ")
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("Analysis : %s", b)
		fmt.Printf("Analysis : %s", b)
	}
}
