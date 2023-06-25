package main

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ClientStats struct {
	FirstSeen int64 `json:"first_seen"`
	LastSeen  int64 `json:"last_seen"`
	Count     int64 `json:"count"`
}

type ClientStatus string

const (
	ClientStatusDefault ClientStatus = "default"
	ClientStatusPending ClientStatus = "pending"
	ClientStatusBlocked ClientStatus = "blocked"
	ClientStatusAllowed ClientStatus = "allowed"
)

type Client struct {
	ID     string       `json:"id"`
	Status ClientStatus `json:"status"`
	Info   *ClientInfo  `json:"info"`
	URL    string       `json:"url"`
	Stats  *ClientStats `json:"stats"`

	update chan bool
}

var (
	mutex   sync.Mutex
	clients map[string]*Client // id to client
)

func init() {
	clients = make(map[string]*Client)
}

func NewClient(ip string, ua string, url string) *Client {
	c := &Client{
		ID:     uuid.New().String(),
		Status: ClientStatusPending,
		URL:    url,
		Stats: &ClientStats{
			FirstSeen: time.Now().UnixMilli(),
			LastSeen:  time.Now().UnixMilli(),
			Count:     1,
		},
		update: make(chan bool, 9),
	}
	log.Println("new client", c.ID)
	c.Info = ParseClientInfo(ip, ua)
	mutex.Lock()
	defer mutex.Unlock()
	clients[c.ID] = c
	return c
}

func GetClientByID(id string) *Client {
	mutex.Lock()
	defer mutex.Unlock()
	if client, ok := clients[id]; ok {
		return client
	}
	return nil
}

func (c *Client) Allow() {
	c.Status = ClientStatusAllowed
	go func() {
		c.update <- true
	}()
}

func (c *Client) Block() {
	c.Status = ClientStatusBlocked
	go func() {
		c.update <- true
	}()
}

func (c *Client) Update() chan bool {
	return c.update
}
