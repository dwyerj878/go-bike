package main

import (
	"bike/analysis"
	"bike/files"
	"bike/models"
	"bike/rider"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	log "github.com/sirupsen/logrus"

	"os"
)

const AuthFailure = "authorization failed"

var currentRide *models.RIDE_DATA
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

	engine := gin.New()
	engine.GET("/chart", chart)
	engine.GET("/data", Authenticate, getData)
	engine.GET("/filename", Authenticate, getFilename)
	engine.POST("/filename", Authenticate, setFilename)
	engine.GET("/datafiles", Authenticate, getFileList)

	engine.Static("/app", "../static/")
	engine.Static("/favicon.ico", "../static/images/favicon.ico")

	log.Info("Starting server")
	showBanner()
	engine.Run(":8081")
}

func showBanner() {
	banner, err := os.ReadFile("banner.txt")
	if err != nil {
		log.Error(err)
	}
	log.Info(string(banner))
}

func Authenticate(context *gin.Context) {
	auth := context.GetHeader("Authorization")
	if auth == "" {
		log.Info("No credentials supplied")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New(AuthFailure)})
		return
	}
	authParts := strings.Split(auth, " ")
	if len(authParts) != 2 {
		log.Info("No credentials supplied")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New(AuthFailure)})
		return
	}
	authMethod, creds := authParts[0], authParts[1]
	if authMethod != "Basic" {
		log.Infof("Auth method %s not supported", authMethod)
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New(AuthFailure)})
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(creds)
	if err != nil {
		log.Error(err)
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New(AuthFailure)})
		return
	}
	if !strings.Contains(string(decoded), ":") {
		log.Info("Invalid credentials format")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New(AuthFailure)})
		return

	}
	credParts := strings.Split(string(decoded), ":")
	username, password := credParts[0], credParts[1]
	if password == allowedUsers[username] {
		context.Set("username", username)
		context.Next()
	} else {
		log.Info("Bad credentials")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New(AuthFailure)})
	}
}

// func getFileList(w http.ResponseWriter, r *http.Request) {
func getFileList(context *gin.Context) {

	filenames, err := files.GetFileList(dataDirectory)
	if err != nil {
		log.Error(err)
	}
	context.JSON(http.StatusOK, filenames)
}

func getFilename(context *gin.Context) {
	log.Debugf("Getting filename %s", fileName)
	context.Writer.Write([]byte("{ \"file_name\" : \"" + fileName + "\" }"))
	return
}

func setFilename(context *gin.Context) {
	r := context.Request

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

func getData(context *gin.Context) {
	b, err := json.MarshalIndent(currentRide.Analysis, "", "  ")
	if err != nil {
		log.Error(err)
	}
	context.Writer.Write(b)
}

func chart(context *gin.Context) {
	w := context.Writer

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
