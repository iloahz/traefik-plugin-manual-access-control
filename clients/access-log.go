package clients

import "time"

type AccessLog struct {
	ClientID  string  `json:"client_id"`
	IPInfo    *IPInfo `json:"ip_info"`
	Host      string  `json:"host"`
	FirstSeen int64   `json:"first_seen"` // milliseconds
	LastSeen  int64   `json:"last_seen"`  // milliseconds
	Count     int64   `json:"count"`
}

func NewAccessLog(clientID string, ip string, host string) *AccessLog {
	return &AccessLog{
		ClientID:  clientID,
		IPInfo:    ParseIPInfo(ip),
		Host:      host,
		FirstSeen: time.Now().UnixMilli(),
		LastSeen:  time.Now().UnixMilli(),
		Count:     1,
	}
}
