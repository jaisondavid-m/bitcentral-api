package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"server/data"
)

type SemesterHandler struct{}

func NewSemesterHandler() *SemesterHandler {
	return &SemesterHandler{}
}

func (h *SemesterHandler) GetSemesterByYear(c *gin.Context) {
	yearParam := c.Param("year")

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid year format",
		})
		return
	}

	semData, exists := data.SemestersData[year]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Year not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"year":    year,
		"count":   len(semData),
		"data":    semData,
	})
}