package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/iloahz/traefik-plugin-manual-access-control/logs"
	"go.uber.org/zap"
)

var (
	r *gin.Engine
)

func CreateServer() {
	r = gin.Default()

	ConfigCORS()
	ServeUI()

	// for traefik plugin
	r.POST("/api/token/generate", generateTokenHandler)
	r.POST("/api/token/validate", validateTokenHandler)

	// for ui
	r.GET("/api/clients", getClientsHandler)
	r.DELETE("/api/client/:id", deleteClientHandler)
	r.GET("/api/client/:id/allow", consentAllowHandler)
	r.GET("/api/client/:id/block", consentBlockHandler)

	server := http.Server{
		Addr:    ":9502",
		Handler: r,
	}

	go func() {
		logs.Info("starting server on port 9502")
		if err := server.ListenAndServe(); err != nil {
			logs.Error("failed to start server", zap.Error(err))
		}
	}()

	// wait SIGINT or SIGTERM and gracefully shutdown the server
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	err := server.Shutdown(context.Background())
	if err != nil {
		logs.Warn("failed to shutdown server", zap.Error(err))
	}
}
