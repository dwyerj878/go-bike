package routes

import (
	"bike/models"
	"bike/rider"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// RoutesTestSuite defines the test suite
type RoutesTestSuite struct {
	suite.Suite
	router       *gin.Engine
	tempDir      string
	testRideFile string
}

// SetupSuite runs once before all tests in the suite
func (s *RoutesTestSuite) SetupSuite() {
	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "test-routes-*")
	s.Require().NoError(err)
	s.tempDir = tempDir
	DataDirectory = tempDir // Set the global DataDirectory for the handlers

	// Create a dummy ride file
	s.testRideFile = "test_ride.json"
	rideContent := `{
		"RIDE": {
			"SAMPLES": [
				{"SECS": 1, "KPH": 10, "WATTS": 100, "CAD": 80},
				{"SECS": 2, "KPH": 20, "WATTS": 200, "CAD": 90}
			],
			"STARTTIME": "2024-01-01T12:00:00Z"
		}
	}`
	err = os.WriteFile(filepath.Join(s.tempDir, s.testRideFile), []byte(rideContent), 0644)
	s.Require().NoError(err)
}

// TearDownSuite runs once after all tests in the suite
func (s *RoutesTestSuite) TearDownSuite() {
	// Clean up the temporary directory
	err := os.RemoveAll(s.tempDir)
	s.Require().NoError(err)
}

// SetupTest runs before each test
func (s *RoutesTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()

	// Register routes without auth middleware for direct handler testing
	s.router.GET("/datafiles", GetFileList)
	s.router.GET("/filename", GetFilename)
	s.router.POST("/filename", SetFilename)
	s.router.GET("/data", GetData)
	s.router.GET("/chart", Chart)

	// Set up initial state for global variables
	FileName = "initial_ride.fit"
	ActiveRider = &rider.RIDER{
		Attributes: []rider.RIDER_ATTRIBUTES{
			{FTP: 250, MaxHR: 190},
		},
	}
	CurrentRide = &models.RideData{
		Ride: models.Ride{
			StartTime: "2023-01-01T10:00:00Z",
			Samples: []models.RideSample{
				{Kph: 25, Watts: 150},
				{Kph: 30, Watts: 250},
			},
		},
		Analysis: models.RideAnalysis{
			AveragePower: 200,
			MaxWatts:     250,
		},
	}
}

// TestRoutesSuite runs the test suite
func TestRoutesSuite(t *testing.T) {
	suite.Run(t, new(RoutesTestSuite))
}

func (s *RoutesTestSuite) TestGetFileList() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/datafiles", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var files []models.FileDetails
	err := json.Unmarshal(w.Body.Bytes(), &files)
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(s.testRideFile, files[0].Filename)
}

func (s *RoutesTestSuite) TestGetFilename() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/filename", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	expectedJSON := fmt.Sprintf(`{ "file_name" : "%s" }`, FileName)
	s.JSONEq(expectedJSON, w.Body.String())
}

func (s *RoutesTestSuite) TestSetFilename() {
	body := strings.NewReader(fmt.Sprintf(`{"file_name": "%s"}`, s.testRideFile))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/filename", body)
	req.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Equal(s.testRideFile, FileName)
	s.NotNil(CurrentRide)
	s.Equal("2024-01-01T12:00:00Z", CurrentRide.Ride.StartTime)
}

func (s *RoutesTestSuite) TestGetData() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/data", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *RoutesTestSuite) TestChart() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/chart", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}
