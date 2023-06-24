package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func waitForConsent(c *gin.Context, client *Client) {
	// check once before waiting
	switch client.Status {
	case ClientStatusAllowed:
		c.JSON(http.StatusOK, gin.H{"token": j.GenerateToken(client.ID)})
		return
	case ClientStatusBlocked:
		c.JSON(http.StatusForbidden, gin.H{"error": "client blocked"})
		return
	}
	for {
		select {
		case <-c.Done():
			return
		case <-client.Update():
			switch client.Status {
			case ClientStatusAllowed:
				c.JSON(http.StatusOK, gin.H{"token": j.GenerateToken(client.ID)})
				return
			case ClientStatusBlocked:
				c.JSON(http.StatusForbidden, gin.H{"error": "client blocked"})
				return
			}
		}
	}
}

func createServer() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	type GenerateTokenRequest struct {
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
	}

	r.POST("/api/token/generate", func(c *gin.Context) {
		var req GenerateTokenRequest
		err := c.Bind(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		client := NewClient(req.IP, req.UserAgent)
		waitForConsent(c, client)
	})

	type ValidateTokenRequest struct {
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
		Token     string `json:"token"`
	}

	r.POST("/api/token/validate", func(c *gin.Context) {
		var req ValidateTokenRequest
		err := c.Bind(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		claims, err := j.ValidateToken(req.Token)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			return
		}
		exp, err := claims.GetExpirationTime()
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			return
		}
		if exp.After(time.Now()) {
			client := GetClientByID(claims.ID)
			if client != nil && client.Info != nil || client.Info.UserAgent == req.UserAgent {
				// normal case for valid token
				// TODO refresh token if exp is close
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
				return
			} else if client == nil {
				// incase token is valid but client is not found, trust it and issue a new token
				client := NewClient(req.IP, req.UserAgent)
				client.Allow()
				waitForConsent(c, client)
				return
			}
		}
		client := NewClient(req.IP, req.UserAgent)
		waitForConsent(c, client)
	})

	r.GET("/api/clients", func(c *gin.Context) {
		var list []*Client
		for _, client := range clients {
			list = append(list, client)
		}
		c.JSON(http.StatusOK, gin.H{"clients": list})
	})

	r.GET("/api/client/:id/allow", func(c *gin.Context) {
		id := c.Param("id")
		client := GetClientByID(id)
		if client == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		client.Allow()
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/api/client/:id/block", func(c *gin.Context) {
		id := c.Param("id")
		client := GetClientByID(id)
		if client == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		client.Block()
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	server := http.Server{
		Addr:    ":9502",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// wait SIGINT or SIGTERM and gracefully shutdown the server
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	err := server.Shutdown(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	createServer()
}
