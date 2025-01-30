package rider

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

func ReadRiderData(fileName string) (*RIDER, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var rider RIDER
	parser := json.NewDecoder(file)
	if err = parser.Decode(&rider); err != nil {
		return nil, err
	}

	createHRZones(&rider)

	createPowerZones(&rider)
	log.Traceln(rider)
	return &rider, nil
}

func createPowerZones(rider *RIDER) {
	ftp := rider.Attributes[0].FTP
	rider.Attributes[0].PowerZones = []RIDER_ZONE{
		{Min: 0, Max: uint32(float64(ftp) * 0.2)},
		{Min: uint32(float64(ftp)*0.2) + 1, Max: uint32(float64(ftp) * 0.5)},
		{Min: uint32(float64(ftp)*0.5) + 1, Max: uint32(float64(ftp) * 0.7)},
		{Min: uint32(float64(ftp)*0.7) + 1, Max: uint32(float64(ftp) * 0.85)},
		{Min: uint32(float64(ftp)*0.85) + 1, Max: uint32(float64(ftp) * 1.0)},
		{Min: uint32(float64(ftp)*1.0) + 1, Max: uint32(float64(ftp) * 1.15)},
		{Min: uint32(float64(ftp)*1.15) + 1, Max: uint32(float64(ftp) * 20.0)},
	}
}

func createHRZones(rider *RIDER) {
	maxHr := rider.Attributes[0].MaxHR
	rider.Attributes[0].HRZones = []RIDER_ZONE{
		{Min: 0, Max: uint32(float64(maxHr) * 0.5)},
		{Min: uint32(float64(maxHr)*0.5) + 1, Max: uint32(float64(maxHr) * 0.6)},
		{Min: uint32(float64(maxHr)*0.6) + 1, Max: uint32(float64(maxHr) * 0.7)},
		{Min: uint32(float64(maxHr)*0.7) + 1, Max: uint32(float64(maxHr) * 0.8)},
		{Min: uint32(float64(maxHr)*0.8) + 1, Max: uint32(float64(maxHr) * 0.9)},
		{Min: uint32(float64(maxHr)*0.9) + 1, Max: uint32(float64(maxHr) * 2.0)},
	}
}
