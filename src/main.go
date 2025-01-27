package main

import (
	"bike/models"
	"encoding/json"
	"fmt"

	"os"
)

func main() {
	println("Hello")
	println(os.Args[1])
	fileName := os.Args[1]
	riderFileName := os.Args[2]
	rider, err := ReadRiderData(riderFileName)
	if err != nil {
		panic(err)
	}
	fmt.Println(rider)
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//var result map[string]interface{}
	var result models.RIDE_DATA
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	//fmt.Println(result.Ride.Tags)
	var max float32
	var powerRanges [100]uint64
	ftp := float32(rider.Attributes[0].FTP)
	wattRange := float32(rider.Attributes[0].FTP / 10)
	overCount := uint64(0)
	underCount := uint64(0)
	zeroCount := uint64(0)
	for _, sample := range result.Ride.Samples {

		if sample.Watts > max {
			max = sample.Watts
		}
		if sample.Watts <= 1.0 {
			zeroCount++
		} else if sample.Watts >= ftp {
			overCount++
		} else {
			underCount++
		}
		idx := int16(sample.Watts / wattRange)
		powerRanges[idx] = powerRanges[idx] + 1

	}
	maxIdx := uint16(max / 25)
	var i uint16
	for i = 0; i < maxIdx; i++ {
		fmt.Printf("range %d %d : count %d\n", i*25, i*25+25, powerRanges[i])

	}

	fmt.Printf("max : %f\n", max)
	fmt.Printf("Zero %d Over %d Under %d\n", zeroCount, overCount, underCount)
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
	fmt.Print(rider.Attributes[0].HRZones)

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
	return &rider, nil
}
