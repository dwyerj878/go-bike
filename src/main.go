package main

import (
	"bike/analysis"
	"bike/files"
	"bike/models"
	"bike/rider"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	log "github.com/sirupsen/logrus"

	"os"
)

var currentRide *models.RideData
var fileName string
var dataDirectory string
var activeRider *rider.RIDER
var allowedUsers map[string]string

func init() {
	allowedUsers = make(map[string]string)
	allowedUsers["admin"] = "password"
	allowedUsers["user"] = "password"
}

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
	// add logging

	http.ListenAndServe(":8081", loggingMiddleware(http.DefaultServeMux))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
		}).Info("http request")
		next.ServeHTTP(w, r)
	})
}

func getImage(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.URL.Path[1:])
	http.ServeFile(w, r, "../static/"+r.URL.Path[1:])
}

func Authenticate(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		log.Info("No credentials supplied")
		return false
	}
	authParts := strings.Split(auth, " ")
	if len(authParts) != 2 {
		log.Info("No credentials supplied")
		return false
	}
	authMethod, creds := authParts[0], authParts[1]
	if authMethod != "Basic" {
		log.Infof("Auth method %s not supported", authMethod)
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(creds)
	if err != nil {
		log.Info("Bad credentials")
		log.Error(err)
		return false
	}
	credParts := strings.Split(string(decoded), ":")
	username, password := credParts[0], credParts[1]
	return password == allowedUsers[username]
}

func getFileList(w http.ResponseWriter, r *http.Request) {
	if !Authenticate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	filenames, err := files.GetFileList(dataDirectory)
	if err != nil {
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"error\" : \"%s\"}", err)))
	}

	out, err := json.MarshalIndent(filenames, "", "  ")
	if err != nil {
		log.Error(err)
		w.Write([]byte(fmt.Sprintf("{\"error\" : \"%s\"}", err)))
	}
	w.Write(out)
}

func getFilename(w http.ResponseWriter, r *http.Request) {
	if !Authenticate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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
		w.Write([]byte("{ \"file_name\" : \"" + fileName + "\" }"))
	}
}

func getData(w http.ResponseWriter, r *http.Request) {
	if !Authenticate(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b, err := json.MarshalIndent(currentRide.Analysis, "", "  ")
	if err != nil {
		log.Error(err)
	}
	w.Write(b)
}

func chart(w http.ResponseWriter, r *http.Request) {

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
