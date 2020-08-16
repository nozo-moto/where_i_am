package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/k0kubun/pp"
	"golang.org/x/xerrors"
)

const (
	ipinfo      = "https://ipinfo.io/ip"
	ipvigilante = "https://ipvigilante.com/"
)

func GetIPAddress() (string, error) {
	body, err := getRequest(ipinfo)
	if err != nil {
		return "", err
	}
	return strings.Replace(string(body), "\n", "", 1), err
}

type LocationRes struct {
	Status string `json:"status"`
	Data   struct {
		Ipv4             string      `json:"ipv4"`
		ContinentName    string      `json:"continent_name"`
		CountryName      string      `json:"country_name"`
		Subdivision1Name interface{} `json:"subdivision_1_name"`
		Subdivision2Name interface{} `json:"subdivision_2_name"`
		CityName         interface{} `json:"city_name"`
		Latitude         string      `json:"latitude"`
		Longitude        string      `json:"longitude"`
	} `json:"data"`
}

func GetLocationFromIP(ipaddress string) (*LocationRes, error) {
	var location LocationRes
	body, err := getRequest(ipvigilante + ipaddress)
	if err != nil {
		return &location, xerrors.Errorf("get_request error %+v", err)
	}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return &location, xerrors.Errorf("json Unmarshal error %+v: body %s", err, string(body))
	}
	return &location, nil
}

func getRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, xerrors.Errorf("request error %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, xerrors.Errorf("response error %w", err)
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func main() {
	ipaddress, err := GetIPAddress()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	location, err := GetLocationFromIP(ipaddress)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	pp.Print(location)
}
