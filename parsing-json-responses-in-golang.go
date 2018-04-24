package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Station struct {
	Id int64 `json:"id"`
	StationName string `json:"stationName"`
	AvailableDocks int64 `json:"availableDocks"`
	TotalDocks int64 `json:"totalDocks"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	StatusValue string `json:"statusValue"`
	StatusKey int64 `json:"statusKey"`
	AvailableBikes int64 `json:"availableBikes"`
	StAddress1 string `json:"stAddress1"`
	StAddress2 string `json:"stAddress2"`
	City string `json:"city"`
	PostalCode string `json:"postalCode"`
	Location string `json:"location"`
	Altitude string `json:"altitude"`
	TestStation bool `json:"testStation"`
	LastCommunicationTime string `json:"lastCommunicationTime"`
	LandMark string `json:"landMark"`
}

type StationAPIResponse struct {
	ExecutionTime string `json:"executionTime"`
	StationBeanList []Station `json:"stationBeanList"`
}

func getStations(body []byte) (*StationAPIResponse, error) {
	var s = new(StationAPIResponse)
	err := json.Unmarshal(body, &s)

	if err != nil {
		fmt.Println("whoops:", err)
	}

	return s, err
}

func main() {
	res, err := http.Get("https://www.citibikenyc.com/stations/json")

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	s, err := getStations([]byte(body))

	//fmt.Println(s.ExecutionTime)
	fmt.Println(s.StationBeanList[0].City)
	fmt.Println(s.StationBeanList[0].StationName)
	fmt.Println(s.StationBeanList[0].Id)
	fmt.Println(s.StationBeanList[0].LandMark)
	fmt.Println(s.StationBeanList[0].Location)
	fmt.Println(s.StationBeanList[0].PostalCode)
	fmt.Println(s.StationBeanList[0].AvailableBikes)
}
