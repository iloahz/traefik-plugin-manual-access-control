package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iloahz/traefik-plugin-manual-access-control/clients"
	"github.com/iloahz/traefik-plugin-manual-access-control/logs"
	"go.uber.org/zap"
)

func updateConsent(allow bool, id, host, ip string) (int, any) {
	logs.Info("update consent", zap.Bool("allow", allow), zap.String("id", id), zap.String("host", host), zap.String("ip", ip))
	client := clients.GetClientByID(id)
	if client == nil {
		logs.Info("client not found")
		return http.StatusNotFound, gin.H{"error": "client not found"}
	}
	if allow {
		client.Allow(host, ip)
	} else {
		client.Block(host, ip)
	}
	return http.StatusOK, gin.H{"status": "ok"}
}

func consentAllowHandler(c *gin.Context) {
	id := c.Param("id")
	host := c.Query("host")
	ip := c.Query("ip")
	if ip == "" {
		ip = clients.AnyIP
	}
	status, resp := updateConsent(true, id, host, ip)
	c.JSON(status, resp)
}

func consentBlockHandler(c *gin.Context) {
	id := c.Param("id")
	host := c.Query("host")
	ip := c.Query("ip")
	if ip == "" {
		ip = clients.AnyIP
	}
	status, resp := updateConsent(false, id, host, ip)
	c.JSON(status, resp)
}
