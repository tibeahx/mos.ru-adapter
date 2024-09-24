package types

type Row struct {
	GlobalID int `json:"global_id"`
	Number   int `json:"Number"`
	Cells    struct {
		Parking Parking
		GeoData GeoData
	} `json:"Cells"`
}

type Parking struct {
	ID                  int    `json:"ID"`
	Name                string `json:"Name"`
	AdmArea             string `json:"AdmArea"`
	GlobalID            string `json:"global_id"`
	District            string `json:"District"`
	Address             string `json:"Address"`
	LocationDescription string `json:"LocationDescription"`
	LongitudeWGS84      string `json:"Longitude_WGS84"`
	LatitudeWGS84       string `json:"Latitude_WGS84"`
	CarCapacity         int    `json:"CarCapacity"`
	Mode                string `json:"Mode"`
}

type GeoData struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}
