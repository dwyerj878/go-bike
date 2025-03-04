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
	Secs   uint64  `json:"SECS"`
	Km     float64 `json:"KM"`
	Watts  float64 `json:"WATTS"`
	Cad    float64 `json:"CAD"`
	Kph    float64 `json:"KPH"`
	Hr     float64 `json:"HR"`
	Alt    float64 `json:"ALT"`
	Lat    float64 `json:"LAT"`
	Long   float64 `json:"LON"`
	Slope  float64 `json:"SLOPE"`
	Temp   float64 `json:"TEMP"`
	Lppb   float64 `json:"LPPB"`
	Rppb   float64 `json:"RPPB"`
	Lppe   float64 `json:"LPPE"`
	Rppe   float64 `json:"RPPE"`
	Lpppb  float64 `json:"LPPPB"`
	Rpppb  float64 `json:"RPPPB"`
	Lpppe  float64 `json:"LPPPE"`
	Rpppe  float64 `json:"RPPPE"`
	Torque float64 `json:"TORQUE"`
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
	MinTemp         float64                       `json:"min_temp"`
	MaxTemp         float64                       `json:"max_temp"`
	MaxWatts        float64                       `json:"max_watts"`
	FTP             RIDE_ANALYSIS_FTP             `json:"ftp"`
	PowerZones      []RIDE_ANALYSIS_ZONE          `json:"power_zones"`
	HRZones         []RIDE_ANALYSIS_ZONE          `json:"hr_zones"`
	NormalizedPower float64                       `json:"normalized_power"`
	AveragePower    float64                       `json:"average_power"`
	AverageSpeed    float64                       `json:"average_speed"`
	ZoneIntervals   []RIDE_ANALYSIS_ZONE_INTERVAL `json:"zone_intervals"`
	AverageCadence  float64                       `json:"average_cadence"`
}

type RIDE_ANALYSIS_FTP struct {
	Over  uint64 `json:"over"`
	Under uint64 `json:"under"`
	Zero  uint64 `json:"zero"`
	FTP   uint32 `json:"white"`
}

type RIDE_ANALYSIS_ZONE struct {
	Min     uint32  `json:"min"`
	Max     uint32  `json:"max"`
	Zone    uint8   `json:"zone"`
	Count   uint64  `json:"count"`
	Percent float64 `json:"percent"`
}

type RIDE_ANALYSIS_ZONE_INTERVAL struct {
	Zone    uint32 `json:"zone"`
	Seconds uint32 `json:"seconds"`
}
