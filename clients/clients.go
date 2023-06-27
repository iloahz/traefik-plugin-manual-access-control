package clients

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID         string       `json:"id"`
	UAInfo     *UAInfo      `json:"ua_info"`
	Consents   []*Consent   `json:"consents"`
	AccessLogs []*AccessLog `json:"access_logs"`

	mutex sync.Mutex
}

var (
	mutex   sync.Mutex
	clients map[string]*Client // id to client
)

func init() {
	clients = make(map[string]*Client)
}

func (c *Client) looksLike(ua string, ip string, host string) bool {
	if c.UAInfo.UserAgent != ua {
		return false
	}
	if accessLog := c.GetAccessLog(ip, host); accessLog == nil {
		return false
	} else {
		return accessLog.IPInfo.IP == ip && time.Now().UnixMilli()-accessLog.LastSeen < 1000*60 // 1 minute
	}
}

func GetClient(ua string, ip string, host string) *Client {
	log.Println("get client", ua, ip, host)
	mutex.Lock()
	defer mutex.Unlock()
	for _, client := range clients {
		if client.looksLike(ua, ip, host) {
			log.Println("existing client", client.ID)
			return client
		}
	}
	client := &Client{
		ID:         uuid.New().String(),
		UAInfo:     ParseUAInfo(ua),
		Consents:   make([]*Consent, 0),
		AccessLogs: make([]*AccessLog, 0),
	}
	client.GetAccessLog(ip, host)
	client.GetConsent(host, AnyIP)
	log.Println("new client", client.ID)
	clients[client.ID] = client
	return client
}

func GetClientByID(id string) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	if client, ok := clients[id]; ok {
		return client
	}
	return nil
}

func DeleteClientByID(id string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(clients, id)
}

func GetAllClients() []*Client {
	var list []*Client
	for _, client := range clients {
		list = append(list, client)
	}
	return list
}

func (c *Client) GetConsent(host string, ip string) *Consent {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, consent := range c.Consents {
		if consent.IP == ip && consent.Host == host {
			return consent
		}
	}
	consent := NewConsent(c.ID, host, ip)
	c.Consents = append(c.Consents, consent)
	return consent
}

func (c *Client) ClearConsents() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Consents = make([]*Consent, 0)
}

func (c *Client) GetAccessLog(ip string, host string) *AccessLog {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, accessLog := range c.AccessLogs {
		if accessLog.IPInfo.IP == ip && accessLog.Host == host {
			return accessLog
		}
	}
	accessLog := NewAccessLog(c.ID, ip, host)
	c.AccessLogs = append(c.AccessLogs, accessLog)
	return accessLog
}

func (c *Client) Allow(host string, ip string) {
	if ip == AnyIP {
		c.ClearConsents()
	}
	c.GetConsent(host, AnyIP).Allow()
}

func (c *Client) Block(host string, ip string) {
	if ip == AnyIP {
		c.ClearConsents()
	}
	c.GetConsent(host, AnyIP).Block()
}
