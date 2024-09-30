package types

import "reflect"

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

func (p Parking) Name() string {
	return p.Cells.Name
}

func (p Parking) Coords() []float64 {
	return p.GeoData.Coordinates
}

type GeoData struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type LocationData struct {
	name      string
	area      string
	district  string
	address   string
	longitude string
	latitude  string
	coords    []float64
}

func (p Parking) Transform() LocationData {
	v := reflect.ValueOf(p.Cells)
	if v.IsNil() {
		return LocationData{}
	}

	coords := p.Coords()

	loc := LocationData{
		name:      p.Cells.Name,
		area:      p.Cells.AmdArea,
		district:  p.Cells.District,
		address:   p.Cells.Address,
		longitude: p.Cells.Longitude_WGS84,
		latitude:  p.Cells.Latitude_WGS84,
		coords:    coords,
	}

	return loc
}
