package routes

import (
	"bike/analysis"
	"bike/files"
	"bike/models"
	"bike/rider"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	log "github.com/sirupsen/logrus"
)

var CurrentRide *models.RideData

var FileName string
var DataDirectory string
var ActiveRider *rider.RIDER

// func getFileList(w http.ResponseWriter, r *http.Request) {
func GetFileList(context *gin.Context) {

	filenames, err := files.GetFileList(DataDirectory)
	if err != nil {
		log.Error(err)
	}
	context.JSON(http.StatusOK, filenames)
}

func GetFilename(context *gin.Context) {
	log.Debugf("Getting filename %s", FileName)
	context.Writer.Write([]byte("{ \"file_name\" : \"" + FileName + "\" }"))
}

func SetFilename(context *gin.Context) {
	r := context.Request

	var request models.LoadRideRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		log.Error(err)
		return
	}
	FileName = request.Filename
	if FileName == "" {
		return
	}
	log.Debugf("Setting filename %s/%s", DataDirectory, FileName)
	fullName := fmt.Sprintf("%s/%s", DataDirectory, FileName)
	ride, err := models.Read(fullName)
	if err != nil {
		panic(err)
	}
	analysis.ExecuteAnalysis(ActiveRider, ride)
	CurrentRide = ride

}

func GetData(context *gin.Context) {
	b, err := json.MarshalIndent(CurrentRide.Analysis, "", "  ")
	if err != nil {
		log.Error(err)
	}
	context.Writer.Write(b)
}

func Chart(context *gin.Context) {
	w := context.Writer

	// create a new line instance
	length := len(CurrentRide.Ride.Samples)
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(

		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic, Width: "1200px", Height: "700px"}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Speed vs Power",
			Subtitle: fmt.Sprintf("Ride data %s", CurrentRide.Ride.StartTime),
		}),
		charts.WithYAxisOpts(
			opts.YAxis{
				SplitNumber: 2,
				Max:         CurrentRide.Analysis.MaxWatts,
				Min:         0,
			}),
		charts.WithXAxisOpts(
			opts.XAxis{
				Max: length,
			}),
	)

	speed := make([]opts.LineData, length)
	power := make([]opts.LineData, length)
	for idx := range length {
		sample := CurrentRide.Ride.Samples[idx]
		speed[idx].Value = sample.Kph
		power[idx].Value = sample.Watts
	}

	// Put data into instance
	line.SetXAxis([]string{"0", "100", "200", "300", "400", "500", "600", "700", "800", "900", "1000"}).
		AddSeries("Power", power).
		AddSeries("Speed", speed).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

	line.Render(w)
}
