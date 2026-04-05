package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"server/data"
)

type LeaveHandler struct{}

func NewLeaveHandler() *LeaveHandler {
	return &LeaveHandler{}
}

func (h *LeaveHandler) GetAllLeaves(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(data.Holidays),
		"data":    data.Holidays,
	})
}