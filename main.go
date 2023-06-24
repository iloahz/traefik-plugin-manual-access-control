package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

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
		for {
			select {
			case <-c.Done():
				return
			case <-client.Update():
				switch client.Status {
				case ClientStatusAllowed:
					c.JSON(http.StatusOK, gin.H{"token": client.GenerateToken()})
					return
				case ClientStatusBlocked:
					c.JSON(http.StatusForbidden, gin.H{"error": "client blocked"})
					return
				}
			}
		}
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
		client := GetClientByToken(req.Token)
		if client == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			return
		}
		if client.Info == nil || client.Info.UserAgent != req.UserAgent {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
