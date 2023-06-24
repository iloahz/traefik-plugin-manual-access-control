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

	r.POST("/api/token/generate", func(c *gin.Context) {
	})

	r.POST("/api/token/validate", func(c *gin.Context) {
	})

	r.GET("/api/clients", func(c *gin.Context) {
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
