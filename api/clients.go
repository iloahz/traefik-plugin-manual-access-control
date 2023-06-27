package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iloahz/traefik-plugin-manual-access-control/clients"
)

func getClientsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"clients": clients.GetAllClients()})
}

func deleteClientHandler(c *gin.Context) {
	id := c.Param("id")
	clients.DeleteClientByID(id)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
