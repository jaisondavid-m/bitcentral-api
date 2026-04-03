package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type RewardActivity struct {
	Date         string `json:"date"`
	RewardPoints string `json:"reward_points"`
	ActivityType string `json:"activity_type"`
	ActivityName string `json:"activity_name"`
}

var rewardSheetTabs = []string{
	"Reward Points Entry",
	"Negative Reward Points",
	"EXTERNAL REWARD POINTS",
}

// Sheet columns (0-indexed):
// A=0 SL.NO | B=1 DATE | C=2 ACTIVITY CODE | D=3 ROLL NO. | E=4 NAME | F=5 YEAR | G=6 DEPT | H=7 REWARD POINTS | I=8 ACTIVITY TYPE | J=9 ACTIVITY NAME

func (h *SheetHandler) GetRewardsByRollNo(c *gin.Context) {
	rawQuery := strings.TrimSpace(c.Query("roll_no"))
	if rawQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Query param 'roll_no' is required",
			"example": "/rewards?roll_no=7376251CS221",
		})
		return
	}

	queryNorm := strings.ToUpper(strings.ReplaceAll(rawQuery, " ", ""))
	spreadsheetID := os.Getenv("RP_Detailed_sheet")
	svc := h.getSheetsService()

	var matched []RewardActivity

	for _, tab := range rewardSheetTabs {
		fullRange := fmt.Sprintf("%s!A1:J100000", tab) // covers up to 100k rows
		resp, err := svc.Spreadsheets.Values.Get(spreadsheetID, fullRange).Do()
		if err != nil {
			continue
		}

		for _, row := range resp.Values {
			rollNorm := strings.ToUpper(strings.ReplaceAll(safeGet(row, 3), " ", ""))
			if rollNorm != queryNorm {
				continue
			}
			matched = append(matched, RewardActivity{
				Date:         safeGet(row, 1),
				RewardPoints: safeGet(row, 7),
				ActivityType: safeGet(row, 8),
				ActivityName: safeGet(row, 9),
			})
		}
	}

	if len(matched) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No reward activities found for this roll number",
		})
		return
	}

	c.JSON(http.StatusOK, matched)
}