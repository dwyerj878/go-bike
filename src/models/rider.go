package models

import "time"

type RIDER struct {
	Name       string             `json:"name"`
	BirthDate  JsonTime           `json:"birthdate"`
	Weight     float32            `json:"weight"`
	Attributes []RIDER_ATTRIBUTES `json:"attributes"`
}

type RIDER_ATTRIBUTES struct {
	FromDate   JsonTime     `json:"from_date"`
	FTP        uint32       `json:"ftp"`
	CP         uint32       `json:"critical_power"`
	MaxHR      uint32       `json:"max_hr"`
	PowerZones []RIDER_ZONE `json:"power_zones,omitempty"`
	HRZones    []RIDER_ZONE `json:"hr_zones,omitempty"`
}

type RIDER_ZONE struct {
	Min uint32 `json:"min"`
	Max uint32 `json:"max"`
}

type JsonTime struct {
	time.Time
}

func (t *JsonTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000-0700"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}
