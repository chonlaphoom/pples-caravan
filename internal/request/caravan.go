package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CaravanInfo struct {
	Url              string
	ResponseStatus   int
	ResponseDuration time.Duration

	Data           CaravanResponse
	VehicleNameMap map[string]string
}

var VehicleNameMap map[string]string

func NewCaravanInfo(url string) *CaravanInfo {
	var vehicleNameMap = map[string]string{
		"67005818": "ลูกน้ำเค็ม,"
		"67006065": "บินหลาดง",
		"67006065": "บินหลาดง",
		"67006066": "นายฮ้อยทมิฬ",
		"67006067": "คมแฝก",
		"67006068": "มนต์รักลูกทุ่ง",
		"67006069": "เพลิงพระนาง",
		"67006070": "กลิ่นกาสะลอง",
	}
	return &CaravanInfo{
		Url:            url,
		VehicleNameMap: vehicleNameMap,
	}
}

func (c *CaravanInfo) String() string {
	s := fmt.Sprintf("Timestamp: %s\n", c.Data.Timestamp)
	for _, v := range c.Data.Data {
		if v.Engine == "ON" {
			v.Engine = "กำลังเดินทาง"
		}
		v.VehicleName = c.VehicleNameMap[v.GpsID]
		s += fmt.Sprintf(`
		Vehicle: %s
			Lat: %.6f|Lon: %.6f
			Speed: %d km/hr|Status: %s|Battery: %sv
			Address: %s
			Last Updated: %s|GPS: %s
			
		`, v.VehicleName, v.Latitude, v.Longitude, v.Speed, v.Engine, v.ExternalBatt, v.AddressT, v.DateTime, v.GPS)
	}
	return strings.ReplaceAll(s, "  ", " ")
}

func (c *CaravanInfo) MakeRequest() (int, time.Duration, error) {
	t := time.Now().Unix() / 10
	url := fmt.Sprintf("%s?t=%d", c.Url, t)

	// log.Printf("Fetching from: %s\n\n", url)
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)

	var result CaravanResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return resp.StatusCode, duration, err
	}

	c.ResponseStatus = resp.StatusCode
	c.ResponseDuration = duration
	c.Data = result

	return resp.StatusCode, duration, nil
}

type CaravanResponse struct {
	Data      []VehicleData `json:"data"`
	Count     int           `json:"count"`
	Filtered  int           `json:"filtered"`
	Total     int           `json:"total"`
	Timestamp string        `json:"timestamp"`
}

type VehicleData struct {
	GpsID            string  `json:"gpsID"`
	PlateNumber      string  `json:"plateNumber"`
	DateTime         string  `json:"dateTime"`
	GPS              string  `json:"GPS"`
	GPRS             string  `json:"GPRS"`
	Engine           string  `json:"Engine"`
	Speed            int     `json:"Speed"`
	Sensor1          string  `json:"Sensor1"`
	Sensor2          string  `json:"Sensor2"`
	Sensor3          string  `json:"Sensor3"`
	Latitude         float64 `json:"Latitude"`
	Longitude        float64 `json:"Longitude"`
	Fuel             int     `json:"Fuel"`
	Temperature      int     `json:"Temperature"`
	COG              int     `json:"COG"`
	VehicleName      string  `json:"vehicleName"`
	VehicleType      string  `json:"vehicleType"`
	GroupVehicle     string  `json:"groupVehicle"`
	IDCard           string  `json:"IDCard"`
	IDTransport      string  `json:"IDTransport"`
	StatusCardReader string  `json:"statusCardReader"`
	Driver           string  `json:"driver"`
	Poi              string  `json:"poi"`
	AddressT         string  `json:"addressT"`
	AddressE         string  `json:"addressE"`
	PowerStatus      string  `json:"powerStatus"`
	ExternalBatt     string  `json:"externalBatt"`
	PositionSource   string  `json:"positionSource"`
}
