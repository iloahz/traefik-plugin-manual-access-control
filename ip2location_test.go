package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetIPInfo(t *testing.T) {
	ip := "116.179.32.218"
	ipInfo, err := GetIPInfo(ip)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	buf, _ := json.Marshal(ipInfo)
	t.Log(string(buf))
}

// {"ip":"116.179.32.218","country_code":"CN","country_name":"China","region_name":"Shanxi","city_name":"Yangquan","latitude":37.8575,"longitude":113.56333,"zip_code":"045000","time_zone":"+08:00","asn":"4837","as":"China169-Backbone China Unicom China169 Backbone","is_proxy":false}
// not using maxmind as it's missing city info for the above test ip
