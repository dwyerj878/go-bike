package main

import (
	"bike/analysis"
	"bike/models"
	"bike/rider"

	"bike/routes"
	"encoding/base64"
	"flag"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"

	"os"
)

const AuthFailure = "authorization failed"

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

	defaultRide := flag.String("defaultRide", "", "Path to the default ride file to load on startup.")
	riderFile := flag.String("riderFile", "", "Path to the rider's JSON data file.")
	dataDirectory := flag.String("dataDirectory", "", "Path to the directory containing ride data files.")
	flag.Parse()

	if *defaultRide == "" || *riderFile == "" || *dataDirectory == "" {
		log.Error("Missing required flags: defaultRide, riderFile, and dataDirectory must be provided.")
		flag.Usage()
		os.Exit(1)
	}

	log.Infof("Loading default ride: %s", *defaultRide)
	routes.FileName = *defaultRide
	routes.DataDirectory = *dataDirectory
	loadedRider, err := rider.ReadRiderData(*riderFile)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	routes.ActiveRider = loadedRider
	ride, err := models.Read(routes.FileName)
	if err != nil {
		panic(err)
	}
	analysis.ExecuteAnalysis(routes.ActiveRider, ride)
	routes.CurrentRide = ride
	log.Debug("http://127.0.0.1:8081/")

	engine := gin.Default()
	CreateRoutes(engine)
	AddStaticContent(engine)

	log.Info("Starting server")
	showBanner()
	engine.Run(":8081")
}

func AddStaticContent(engine *gin.Engine) {
	engine.Static("/app", "../static/")
	engine.Static("/favicon.ico", "../static/images/favicon.ico")
}

func CreateRoutes(engine *gin.Engine) {
	engine.GET("/chart", routes.Chart)
	engine.GET("/data", Authenticate, routes.GetData)
	engine.GET("/filename", Authenticate, routes.GetFilename)
	engine.POST("/filename", Authenticate, routes.SetFilename)
	engine.GET("/datafiles", Authenticate, routes.GetFileList)
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
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": AuthFailure})
		return
	}
	authParts := strings.Split(auth, " ")
	if len(authParts) != 2 {
		log.Info("No credentials supplied")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": AuthFailure})
		return
	}
	authMethod, creds := authParts[0], authParts[1]
	if authMethod != "Basic" {
		log.Infof("Auth method %s not supported", authMethod)
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": AuthFailure})
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(creds)
	if err != nil {
		log.Error(err)
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": AuthFailure})
		return
	}
	if !strings.Contains(string(decoded), ":") {
		log.Info("Invalid credentials format")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": AuthFailure})
		return

	}
	credParts := strings.Split(string(decoded), ":")
	username, password := credParts[0], credParts[1]
	if password == allowedUsers[username] {
		context.Set("username", username)
		context.Next()
	} else {
		log.Info("Bad credentials")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": AuthFailure})
	}
}
