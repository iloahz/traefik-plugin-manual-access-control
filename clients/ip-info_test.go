package clients

import (
	"encoding/json"
	"testing"
)

func TestParseIPInfo(t *testing.T) {
	ip := "116.179.32.218"
	info := ParseIPInfo(ip)
	buf, _ := json.MarshalIndent(info, "", "  ")
	t.Log(string(buf))
}

// {
// 	"ip": "116.179.32.218",
// 	"country": "China",
// 	"region": "Shanxi",
// 	"city": "Yangquan",
// 	"lat": 37.8575,
// 	"long": 113.56333,
// 	"as": "China169-Backbone China Unicom China169 Backbone",
// 	"asn": "4837"
// }
