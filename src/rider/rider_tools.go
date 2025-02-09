package rider

import (
	"encoding/json"
	"math"
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
		{Min: 0, Max: calcPercent(ftp, 0.2)},
		{Min: calcPercent(ftp, 0.2) + 1, Max: calcPercent(ftp, 0.5)},
		{Min: calcPercent(ftp, 0.5) + 1, Max: calcPercent(ftp, 0.7)},
		{Min: calcPercent(ftp, 0.7) + 1, Max: calcPercent(ftp, 0.85)},
		{Min: calcPercent(ftp, 0.85) + 1, Max: ftp},
		{Min: ftp + 1, Max: calcPercent(ftp, 1.15)},
		{Min: calcPercent(ftp, 1.15) + 1, Max: calcPercent(ftp, 20.0)},
	}
}

func calcPercent(base uint32, multiplier float64) uint32 {
	val := math.Round(float64(base) * multiplier)
	return uint32(val)
}

func createHRZones(rider *RIDER) {
	maxHr := rider.Attributes[0].MaxHR
	rider.Attributes[0].HRZones = []RIDER_ZONE{
		{Min: 0, Max: calcPercent(maxHr, 0.5)},
		{Min: calcPercent(maxHr, 0.5) + 1, Max: calcPercent(maxHr, 0.6)},
		{Min: calcPercent(maxHr, 0.6) + 1, Max: calcPercent(maxHr, 0.7)},
		{Min: calcPercent(maxHr, 0.7) + 1, Max: calcPercent(maxHr, 0.8)},
		{Min: calcPercent(maxHr, 0.8) + 1, Max: calcPercent(maxHr, 0.9)},
		{Min: calcPercent(maxHr, 0.9) + 1, Max: calcPercent(maxHr, 2.0)},
	}
}
