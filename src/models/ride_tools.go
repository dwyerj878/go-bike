package models

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

/*
Read the ride data from a json file
*/
func Read(fileName string) (*RIDE_DATA, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//var ride map[string]interface{}
	var ride RIDE_DATA
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&ride)
	if err != nil {
		return nil, err
	}
	log.Traceln(ride)
	return &ride, nil

}
