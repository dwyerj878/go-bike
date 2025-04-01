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
var fileName string
var dataDirectory string
var activeRider *rider.RIDER

func main() {

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(
		&log.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
			PadLevelText:  true,
		})
	log.SetReportCaller(true)

	log.Info(os.Args[1])
	fileName = os.Args[1]
	riderFileName := os.Args[2]
	dataDirectory = os.Args[3]
	loadedRider, err := rider.ReadRiderData(riderFileName)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	activeRider = loadedRider
	ride, err := models.Read(fileName)
	if err != nil {
		panic(err)
	}
	analysis.ExecuteAnalysis(activeRider, ride)
	currentRide = ride
	log.Debug("http://127.0.0.1:8081/")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "../static/index.html") })
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "../static/style.css") })
	http.HandleFunc("/images/", getImage)
	http.HandleFunc("/favicon.ico", getImage)
	http.HandleFunc("/chart", chart)
	http.HandleFunc("/data", getData)
	http.HandleFunc("/filename", getFilename)
	http.HandleFunc("/datafiles", getFileList)
	log.Info("Starting server")
	http.ListenAndServe(":8081", nil)

}

func getImage(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.URL.Path[1:])
	http.ServeFile(w, r, "../static/"+r.URL.Path[1:])
}

func getFileList(w http.ResponseWriter, r *http.Request) {
	dirEntries, err := os.ReadDir(dataDirectory)
	if err != nil {
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"error\" : \"%s\"}", err)))
	}
	filenames := []string{}
	for _, entry := range dirEntries {
		filenames = append(filenames, entry.Name())
	}
	out, err := json.MarshalIndent(filenames, "", "  ")
	if err != nil {
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"error\" : \"%s\"}", err)))
	}
	w.Write(out)
}

func getFilename(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Debugf("Getting filename %s", fileName)
		w.Write([]byte("{ \"file_name\" : \"" + fileName + "\" }"))
		return
	}

	if r.Method == http.MethodPost {
		var request models.LoadRideRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			log.Error(err)
			return
		}
		fileName = request.Filename
		if fileName == "" {
			return
		}
		log.Debugf("Setting filename %s/%s", dataDirectory, fileName)
		fullName := fmt.Sprintf("%s/%s", dataDirectory, fileName)
		ride, err := models.Read(fullName)
		if err != nil {
			panic(err)
		}
		analysis.ExecuteAnalysis(activeRider, ride)
		currentRide = ride
	}
}

func getData(w http.ResponseWriter, r *http.Request) {
	b, err := json.MarshalIndent(currentRide.Analysis, "", "  ")
	if err != nil {
		log.Error(err)
	}
	w.Write(b)
}

func chart(w http.ResponseWriter, _ *http.Request) {
	// create a new line instance
	length := len(currentRide.Ride.Samples)
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(

		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic, Width: "1200px", Height: "700px"}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Speed vs Power",
			Subtitle: fmt.Sprintf("Ride data %s", currentRide.Ride.StartTime),
		}),
		charts.WithYAxisOpts(
			opts.YAxis{
				SplitNumber: 2,
				Max:         currentRide.Analysis.MaxWatts,
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
		sample := currentRide.Ride.Samples[idx]
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
