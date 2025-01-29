package models

type RIDE_DATA struct {
	Ride     RIDE          `json:"RIDE"`
	Analysis RIDE_ANALYSIS `json:"Analysis"`
}

type RIDE struct {
	StartTime  string          `json:"STARTTIME"`
	RecIntSecs int             `json:"RECINTSECS"`
	DeviceType string          `json:"DEVICETYPE"`
	Identifier string          `json:"IDENTIFIER"`
	Tags       RIDE_TAGS       `json:"TAGS"`
	Intervals  []RIDE_INTERVAL `json:"INTERVALS"`
	Samples    []RIDE_SAMPLE   `json:"SAMPLES"`
}

type RIDE_SAMPLE struct {
	Secs  uint64  `json:"SECS"`
	Km    float32 `json:"KM"`
	Watts float32 `json:"WATTS"`
	Cad   float32 `json:"CAD"`
	Kph   float32 `json:"KPH"`
	Hr    float32 `json:"HR"`
	Alt   float32 `json:"ALT"`
	Lat   float32 `json:"LAT"`
	Long  float32 `json:"LON"`
	Slope float32 `json:"SLOPE"`
	Temp  float32 `json:"TEMP"`
	Lppb  float32 `json:"LPPB"`
	Rppb  float32 `json:"RPPB"`
	Lppe  float32 `json:"LPPE"`
	Rppe  float32 `json:"RPPE"`
	Lpppb float32 `json:"LPPPB"`
	Rpppb float32 `json:"RPPPB"`
	Lpppe float32 `json:"LPPPE"`
	Rpppe float32 `json:"RPPPE"`
}

type RIDE_INTERVAL struct {
	Name  string `json:"NAME"`
	Start int    `json:"START"`
	Stop  int    `json:"STOP"`
	Color string `json:"COLOR"`
	Ptest string `json:"PTEST"`
}

type RIDE_TAGS struct {
	Athlete                 string `json:"Athlete"`
	AnaerobicTrainingEffect string `json:"Anaerobic Training Effect"`

	AerobicTrainingEffect string `json:"Aerobic Training Effect"`
	ChangeHistory         string `json:"Change History"`
	Data                  string `json:"Data"`
	Device                string `json:"Device"`
	DeviceInfo            string `json:"Device Info"`
}

type RIDE_ANALYSIS struct {
	MinTemp  float32
	MaxTemp  float32
	MaxWatts float32
	FTP      RIDE_ANALYSIS_FTP
	ZONES    []RIDE_ANALYSIS_ZONE
}

type RIDE_ANALYSIS_FTP struct {
	Over  uint64
	Under uint64
	Zero  uint64
	FTP   uint32
}

type RIDE_ANALYSIS_ZONE struct {
	Zone  uint8
	Count uint64
}
