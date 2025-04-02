package models

type LoadRideRequest struct {
	Filename string `json:"file_name"`
}

type FileDetails struct {
	Filename   string `json:"file_name"`
	ModifyDate string `json:"modify_date"`
	Size       int64  `json:"size"`
}
