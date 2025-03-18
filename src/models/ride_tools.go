package models

import (
	"encoding/json"
	"fmt"
	"os"

	"path/filepath"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/untyped/fieldnum"
	log "github.com/sirupsen/logrus"
)

/*
Read the ride data from a json file
*/
func Read(fileName string) (*RIDE_DATA, error) {
	if fileName == "" {
		return nil, nil
	}
	extension := filepath.Ext(fileName)

	if extension == ".json" {
		return readJsonFile(fileName)
	}

	if extension == ".fit" {
		return readFitFile(fileName)
	}

	return nil, fmt.Errorf("invalid file type %s", fileName)
}

func readJsonFile(fileName string) (*RIDE_DATA, error) {
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

func readFitFile(fileName string) (*RIDE_DATA, error) {
	ride := RIDE_DATA{
		Ride: RIDE{
			Samples: []RIDE_SAMPLE{},
		},
	}
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s %s", fileName, err.Error())
	}
	defer f.Close()

	dec := decoder.New(f)

	fit, err := dec.Decode()
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s %s", fileName, err.Error())
	}

	fmt.Printf("FileHeader DataSize: %d\n", fit.FileHeader.DataSize)
	fmt.Printf("Messages count: %d\n", len(fit.Messages))
	// FileId is always the first message; 4 = activity
	fmt.Printf("File Type: %v\n", fit.Messages[0].FieldValueByNum(fieldnum.FileIdType).Any())

	activity := filedef.NewActivity(fit.Messages...)

	fmt.Printf("Sessions count: %d\n", len(activity.Sessions))
	fmt.Printf("Laps count: %d\n", len(activity.Laps))
	fmt.Printf("Records count: %d\n", len(activity.Records))

	for _, record := range activity.Records {
		sample := RIDE_SAMPLE{
			Secs:   uint64(record.TimestampUint32()),
			Km:     record.DistanceScaled(),
			Watts:  float64(record.Power),
			Cad:    float64(record.Cadence),
			Kph:    record.SpeedScaled() * 3.6,
			Hr:     float64(record.HeartRate),
			Alt:    record.AltitudeScaled(),
			Lat:    record.PositionLatDegrees(),
			Long:   record.PositionLongDegrees(),
			Slope:  record.GradeScaled(),
			Temp:   float64(record.Temperature),
			Lppb:   0,
			Rppb:   0,
			Lppe:   0,
			Rppe:   0,
			Lpppb:  0,
			Rpppb:  0,
			Lpppe:  0,
			Rpppe:  0,
			Torque: 0,
		}
		if sample.Watts > 65530 {
			sample.Watts = 0
		}

		ride.Ride.Samples = append(ride.Ride.Samples, sample)

	}

	return &ride, nil
}
