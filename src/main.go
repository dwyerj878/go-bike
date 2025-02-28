package main

import (
	"bike/analysis"
	"bike/models"
	"bike/rider"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	log "github.com/sirupsen/logrus"

	"os"
)

var currentRide *models.RIDE_DATA

func main() {

	log.SetLevel(log.DebugLevel)

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
	currentRide = ride

	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	// create a new line instance
	length := len(currentRide.Ride.Samples)
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line example in Westeros theme",
			Subtitle: "Line chart rendered by the http server this time",
		}),
		charts.WithYAxisOpts(
			opts.YAxis{
				SplitNumber: 2,
				Max:         currentRide.Analysis.MaxWatts,
				Min:         0,
			}), charts.WithXAxisOpts(
			opts.XAxis{
				Max: length,
			}))

	speed := make([]opts.LineData, length)
	power := make([]opts.LineData, length)
	for idx := range length {
		sample := currentRide.Ride.Samples[idx]
		speed[idx].Value = sample.Kph
		power[idx].Value = sample.Watts
	}

	// Put data into instance
	line.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Power", power).
		AddSeries("Speed", speed).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

	line.Render(w)
}
