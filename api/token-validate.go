package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iloahz/traefik-plugin-manual-access-control/clients"
	"github.com/iloahz/traefik-plugin-manual-access-control/logs"
	"github.com/iloahz/traefik-plugin-manual-access-control/tokens"
	"go.uber.org/zap"
)

type ValidateTokenRequest struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Host      string `json:"host"`
	Token     string `json:"token"`
}

type ValidateTokenResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func blockWithNewToken(consent *clients.Consent) (int, any) {
	return http.StatusForbidden, ValidateTokenResponse{
		ID:    consent.ClientID,
		Token: consent.GenerateToken(),
	}
}

func validateToken(req *ValidateTokenRequest) (int, any) {
	logs.Info("validate token", zap.Any("req", req))
	claims, err := tokens.ValidateToken(req.Token)
	if err != nil {
		logs.Info("invalid token, error:", zap.Error(err))
		client := clients.GetClient(req.IP, req.UserAgent, req.Host)
		return blockWithNewToken(client.GetConsent(req.Host, clients.AnyIP))
	}
	client := clients.GetClientByID(claims.ID)
	if client == nil {
		// token is valid but client is not found, trust the claims
		// TODO disallow this case after supporting persistent storage
		fmt.Println("token is valid but client is not found")
		client = clients.GetClient(req.IP, req.UserAgent, req.Host)
		consent := client.GetConsent(claims.Host, claims.IP)
		consent.Allow()
	}
	exp, err := claims.GetExpirationTime()
	if err != nil {
		logs.Info("error getting expiration time, error:", zap.Error(err))
		return blockWithNewToken(client.GetConsent(req.Host, clients.AnyIP))
	}
	if exp.Before(time.Now()) {
		logs.Info("token expired")
		client.ClearConsents()
		return blockWithNewToken(client.GetConsent(req.Host, clients.AnyIP))
	}
	// up to here, token is valid and not expired
	if client.UAInfo == nil || client.UAInfo.UserAgent != req.UserAgent {
		logs.Info("token is valid but client info does not match")
		return blockWithNewToken(client.GetConsent(req.Host, clients.AnyIP))
	}
	consentIP := client.GetConsent(req.Host, req.IP)
	if consentIP.Status == clients.ConsentStatusAllowed {
		// token is valid and client info matches and consent is allowed
		// TODO refresh token if exp is close
		logs.Info("token is valid and client info matches and consent is allowed")
		return http.StatusOK, ValidateTokenResponse{
			ID:    client.ID,
			Token: req.Token,
		}
	} else if consentIP.Status == clients.ConsentStatusBlocked {
		return blockWithNewToken(client.GetConsent(req.Host, req.IP))
	}
	consent := client.GetConsent(req.Host, clients.AnyIP)
	if consent.Status == clients.ConsentStatusAllowed {
		// token is valid and client info matches and consent is allowed
		// TODO refresh token if exp is close
		logs.Info("token is valid and client info matches and consent is allowed")
		return http.StatusOK, ValidateTokenResponse{
			ID:    client.ID,
			Token: req.Token,
		}
	} else if consent.Status == clients.ConsentStatusBlocked {
		return blockWithNewToken(client.GetConsent(req.Host, clients.AnyIP))
	}
	// consent pending
	fmt.Println("token is valid but consent is pending")
	return blockWithNewToken(client.GetConsent(req.Host, clients.AnyIP))
}

func validateTokenHandler(c *gin.Context) {
	var req ValidateTokenRequest
	err := c.Bind(&req)
	if err != nil {
		logs.Info("error parsing request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	status, resp := validateToken(&req)
	c.JSON(status, resp)
}
