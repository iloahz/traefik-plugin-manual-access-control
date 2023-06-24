package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type IPInfo struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionName  string  `json:"region_name"`
	CityName    string  `json:"city_name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	ASN         string  `json:"asn"`
	AS          string  `json:"as"`
	IsProxy     bool    `json:"is_proxy"`
}

func GetIPInfo(ip string) (*IPInfo, error) {
	apiKey := os.Getenv("IP2LOCATION_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("IP2LOCATION_API_KEY environment variable is not set")
	}

	url := fmt.Sprintf("https://api.ip2location.io/?key=%s&ip=%s", apiKey, ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ipInfo IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return nil, err
	}

	return &ipInfo, nil
}
