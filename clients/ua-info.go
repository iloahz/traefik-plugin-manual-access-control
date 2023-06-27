package clients

import "github.com/mileusna/useragent"

type UAInfo struct {
	UserAgent      string `json:"user_agent"`
	BrowserName    string `json:"browser_name"`
	BrowserVersion string `json:"browser_version"`
	OSName         string `json:"os_name"`
	OSVersion      string `json:"os_version"`
}

func ParseUAInfo(userAgent string) *UAInfo {
	info := &UAInfo{
		UserAgent: userAgent,
	}

	ua := useragent.Parse(userAgent)
	info.BrowserName = ua.Name
	info.BrowserVersion = ua.Version
	info.OSName = ua.OS
	info.OSVersion = ua.OSVersion

	return info
}
