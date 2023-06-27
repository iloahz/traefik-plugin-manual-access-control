package clients

import (
	"encoding/json"
	"testing"
)

func TestParseUAInfo(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	info := ParseUAInfo(ua)
	buf, _ := json.MarshalIndent(info, "", "  ")
	t.Log(string(buf))
}

// {
// 	"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
// 	"browser_name": "Chrome",
// 	"browser_version": "114.0.0.0",
// 	"os_name": "macOS",
// 	"os_version": "10.15.7"
// }
