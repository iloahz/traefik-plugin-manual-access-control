package main

import (
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	Stats  *ClientStats `json:"stats"`

	update chan bool
}

var (
	mutex   sync.Mutex
	clients map[string]*Client // id to client
)

var (
	j *JWT
)

func init() {
	var err error
	j, err = NewJWT(os.Getenv("JWT_SECRET"))
	if err != nil {
		panic(err)
	}
	clients = make(map[string]*Client)
}

func NewClient(ip string, ua string) *Client {
	c := &Client{
		ID:     uuid.New().String(),
		Status: ClientStatusPending,
		Info:   ParseClientInfo(ip, ua),
		Stats: &ClientStats{
			FirstSeen: time.Now().UnixMilli(),
			LastSeen:  time.Now().UnixMilli(),
			Count:     1,
		},
	}
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

func GetClientByToken(token string) *Client {
	if len(token) > 0 {
		if claims, err := j.ValidateToken(token); err == nil {
			if mapClaims, ok := claims.(jwt.MapClaims); ok {
				if id, ok := mapClaims["id"].(string); ok {
					return GetClientByID(id)
				}
			}
		}
	}
	// TODO deal with the case that the token is valid but the client is not found
	return nil
}

func (c *Client) Allow() {
	c.Status = ClientStatusAllowed
	c.update <- true
}

func (c *Client) Block() {
	c.Status = ClientStatusBlocked
	c.update <- true
}

func (c *Client) Update() chan bool {
	return c.update
}

func (c *Client) GenerateToken() string {
	return j.GenerateToken(c.ID)
}
