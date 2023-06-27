package clients

type IPInfo struct {
	IP      string  `json:"ip"`
	Country string  `json:"country"`
	Region  string  `json:"region"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
	AS      string  `json:"as"`
	ASN     string  `json:"asn"`
}

func ParseIPInfo(ip string) *IPInfo {
	info := &IPInfo{
		IP: ip,
	}

	if ipInfo, err := GetIPInfo(ip); err == nil {
		info.Country = ipInfo.CountryName
		info.Region = ipInfo.RegionName
		info.City = ipInfo.CityName
		info.Lat = ipInfo.Latitude
		info.Long = ipInfo.Longitude
		info.AS = ipInfo.AS
		info.ASN = ipInfo.ASN
	}

	return info
}
