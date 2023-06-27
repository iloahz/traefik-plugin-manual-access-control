package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iloahz/traefik-plugin-manual-access-control/clients"
)

type GenerateTokenRequest struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Host      string `json:"host"`
}

type GenerateTokenResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func generateTokenHandler(c *gin.Context) {
	var req GenerateTokenRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client := clients.GetClient(req.IP, req.UserAgent, req.Host)
	consent := client.GetConsent(req.Host, clients.AnyIP)
	c.JSON(http.StatusOK, GenerateTokenResponse{
		ID:    client.ID,
		Token: consent.GenerateToken(),
	})
}
