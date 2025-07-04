package models

type RideData struct {
	Ride     Ride         `json:"RIDE"`
	Analysis RideAnalysis `json:"Analysis"`
}

type Ride struct {
	StartTime  string         `json:"STARTTIME"`
	RecIntSecs int            `json:"RECINTSECS"`
	DeviceType string         `json:"DEVICETYPE"`
	Identifier string         `json:"IDENTIFIER"`
	Tags       RideTags       `json:"TAGS"`
	Intervals  []RideInterval `json:"INTERVALS"`
	Samples    []RideSample   `json:"SAMPLES"`
}

type RideSample struct {
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

type RideInterval struct {
	Name  string `json:"NAME"`
	Start int    `json:"START"`
	Stop  int    `json:"STOP"`
	Color string `json:"COLOR"`
	Ptest string `json:"PTEST"`
}

type RideTags struct {
	Athlete                 string `json:"Athlete"`
	AnaerobicTrainingEffect string `json:"Anaerobic Training Effect"`

	AerobicTrainingEffect string `json:"Aerobic Training Effect"`
	ChangeHistory         string `json:"Change History"`
	Data                  string `json:"Data"`
	Device                string `json:"Device"`
	DeviceInfo            string `json:"Device Info"`
}

type RideAnalysis struct {
	MinTemp         float64                    `json:"min_temp"`
	MaxTemp         float64                    `json:"max_temp"`
	MaxWatts        float64                    `json:"max_watts"`
	FTP             RideAnalysisFtp            `json:"ftp"`
	PowerZones      []RideAnalysisZone         `json:"power_zones"`
	HRZones         []RideAnalysisZone         `json:"hr_zones"`
	NormalizedPower float64                    `json:"normalized_power"`
	AveragePower    float64                    `json:"average_power"`
	AverageSpeed    float64                    `json:"average_speed"`
	ZoneIntervals   []RideAnalysisZoneInterval `json:"zone_intervals"`
	AverageCadence  float64                    `json:"average_cadence"`

	TSS float64 `json:"tss"`
	IFF float64 `json:"iff"`
}

type RideAnalysisFtp struct {
	Over  uint64 `json:"over"`
	Under uint64 `json:"under"`
	Zero  uint64 `json:"zero"`
	FTP   uint32 `json:"ftp"`
}

type RideAnalysisZone struct {
	Min     uint32  `json:"min"`
	Max     uint32  `json:"max"`
	Zone    uint8   `json:"zone"`
	Count   uint64  `json:"count"`
	Percent float64 `json:"percent"`
}

type RideAnalysisZoneInterval struct {
	Zone    uint32 `json:"zone"`
	Seconds uint32 `json:"seconds"`
}
