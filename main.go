package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func createServer() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(static.Serve("/", static.LocalFile("ui/dist", false)))

	type GenerateTokenRequest struct {
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
		URL       string `json:"url"`
	}

	type GenerateTokenResponse struct {
		ID    string `json:"id"`
		Token string `json:"token"`
	}

	r.POST("/api/token/generate", func(c *gin.Context) {
		var req GenerateTokenRequest
		err := c.Bind(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		client := NewClient(req.IP, req.UserAgent, req.URL)
		c.JSON(http.StatusOK, GenerateTokenResponse{
			ID:    client.ID,
			Token: j.GenerateToken(client.ID),
		})
	})

	type ValidateTokenRequest struct {
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
		URL       string `json:"url"`
		Token     string `json:"token"`
	}

	type ValidateTokenResponse struct {
		ID    string `json:"id"`
		Token string `json:"token"`
	}

	r.POST("/api/token/validate", func(c *gin.Context) {
		fmt.Println("validate token")
		var req ValidateTokenRequest
		err := c.Bind(&req)
		if err != nil {
			fmt.Println("error parsing request", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		fmt.Println("req", req)
		claims, err := j.ValidateToken(req.Token)
		if err != nil {
			fmt.Println("invalid token, error:", err)
			client := NewClient(req.IP, req.UserAgent, req.URL)
			c.JSON(http.StatusForbidden, ValidateTokenResponse{
				ID:    client.ID,
				Token: j.GenerateToken(client.ID),
			})
			return
		}
		exp, err := claims.GetExpirationTime()
		if err != nil {
			fmt.Println("error getting expiration time, error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting expiration time"})
			return
		}
		if exp.Before(time.Now()) {
			fmt.Println("token expired")
			client := GetClientByID(claims.ID)
			if client == nil {
				client = NewClient(req.IP, req.UserAgent, req.URL)
			}
			c.JSON(http.StatusForbidden, ValidateTokenResponse{
				ID:    client.ID,
				Token: j.GenerateToken(client.ID),
			})
			return
		}
		client := GetClientByID(claims.ID)
		if client == nil {
			fmt.Println("token is valid but client is not found")
			// token is valid but client is not found, trust it and issue a new token
			// TODO disallow this case after supporting persistent storage
			client = NewClient(req.IP, req.UserAgent, req.URL)
			client.Allow()
			c.JSON(http.StatusOK, ValidateTokenResponse{
				ID:    client.ID,
				Token: j.GenerateToken(client.ID),
			})
			return
		}
		if client.Info != nil && client.Info.UserAgent == req.UserAgent && client.Status == ClientStatusAllowed {
			fmt.Println("token is valid and client info matches and client is allowed")
			// normal case for valid token
			// TODO refresh token if exp is close
			c.JSON(http.StatusOK, ValidateTokenResponse{
				ID:    client.ID,
				Token: req.Token,
			})
			return
		}
		client = NewClient(req.IP, req.UserAgent, req.URL)
		c.JSON(http.StatusForbidden, ValidateTokenResponse{
			ID:    client.ID,
			Token: j.GenerateToken(client.ID),
		})
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
		fmt.Println("starting server")
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
	if os.Getenv("DEBUG") == "true" {
		go func() {
			time.Sleep(time.Second)
			NewClient("116.179.32.218", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36", "cdn.home.iloahz.com")
			c2 := NewClient("221.192.199.49", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1", "code.home.iloahz.com")
			c2.Allow()
			c3 := NewClient("180.163.220.66", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36", "chatgpt.home.iloahz.com")
			c3.Block()
		}()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	createServer()
}
