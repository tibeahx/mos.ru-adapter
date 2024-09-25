package types

type Metadata struct {
	GlobalId int `json:"global_id"`
	Number   int `json:"Number"`
}

type Cells struct {
	ID                  int    `json:"ID"`
	Name                string `json:"Name"`
	GlobalId            int    `json:"global_id"`
	AmdArea             string `json:"AmdArea"`
	District            string `json:"District"`
	Address             string `json:"Address"`
	LocationDescription string `json:"LocationDescription"`
	Longitude_WGS84     string `json:"Longitude_WGS84"`
	Latitude_WGS84      string `json:"Latitude_WGS84"`
	CarCapacity         int    `json:"CarCapacity"`
	Mode                string `json:"Mode"`
}

type Parking struct {
	Metadata
	Cells   Cells   `json:"Cells"`
	GeoData GeoData `json:"geoData"`
}

type GeoData struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}
