package handlers

import (
	"net/http"

	"server/data"

	"github.com/gin-gonic/gin"
)

type CardHandler struct{}

func NewCardHandler() *CardHandler {
	return &CardHandler{}
}

func GetCards(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(data.Cards),
		"data":    data.Cards,
	})
}