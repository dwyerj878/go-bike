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
	ftp := float32(250.0)
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
		idx := int16(sample.Watts / 25)
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
