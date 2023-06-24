package main

import (
	"encoding/json"
	"testing"
)

func TestParseClientInfo(t *testing.T) {
	ip := "116.179.32.218"
	clientInfo := ParseClientInfo(ip, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	// print clientInfo as pretty json
	buf, _ := json.MarshalIndent(clientInfo, "", "  ")
	t.Log(string(buf))
}

// {
// 	"ip": "116.179.32.218",
// 	"country": "China",
// 	"region": "Shanxi",
// 	"city": "Yangquan",
// 	"asn": "China169-Backbone China Unicom China169 Backbone",
// 	"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
// 	"browser_name": "Chrome",
// 	"browser_version": "114.0.0.0",
// 	"os_name": "macOS",
// 	"os_version": "10.15.7"
// }
