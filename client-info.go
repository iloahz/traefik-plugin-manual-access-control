package main

import (
	"github.com/mileusna/useragent"
)

type ClientInfo struct {
	IP      string  `json:"ip"`
	Country string  `json:"country"`
	Region  string  `json:"region"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
	ASN     string  `json:"asn"`

	UserAgent      string `json:"user_agent"`
	BrowserName    string `json:"browser_name"`
	BrowserVersion string `json:"browser_version"`
	OSName         string `json:"os_name"`
	OSVersion      string `json:"os_version"`
}

func ParseClientInfo(ip string, userAgent string) *ClientInfo {
	c := &ClientInfo{
		IP:        ip,
		UserAgent: userAgent,
	}

	if ipInfo, err := GetIPInfo(ip); err == nil {
		c.Country = ipInfo.CountryName
		c.Region = ipInfo.RegionName
		c.City = ipInfo.CityName
		c.Lat = ipInfo.Latitude
		c.Long = ipInfo.Longitude
		c.ASN = ipInfo.AS
	}

	ua := useragent.Parse(userAgent)
	c.BrowserName = ua.Name
	c.BrowserVersion = ua.Version
	c.OSName = ua.OS
	c.OSVersion = ua.OSVersion

	return c
}
