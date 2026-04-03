package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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

const rewardsCacheTTL = 5 * time.Minute

var (
	rewardsCacheMu      sync.RWMutex
	rewardsCacheBuildMu sync.Mutex
	rewardsByRollCache  map[string][]RewardActivity
	rewardsCacheExpiry  time.Time
)

func normalizeRollNo(s string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(s), " ", ""))
}

func cloneRewardMap(src map[string][]RewardActivity) map[string][]RewardActivity {
	if src == nil {
		return nil
	}
	out := make(map[string][]RewardActivity, len(src))
	for k, v := range src {
		vv := make([]RewardActivity, len(v))
		copy(vv, v)
		out[k] = vv
	}
	return out
}

func (h *SheetHandler) buildRewardsIndex() (map[string][]RewardActivity, error) {
	spreadsheetID := os.Getenv("RP_Detailed_sheet")
	svc := h.getSheetsService()

	if svc == nil {
		return nil, fmt.Errorf("sheets service not initialized")
	}

	ranges := make([]string, 0, len(rewardSheetTabs))
	for _, tab := range rewardSheetTabs {
		ranges = append(ranges, fmt.Sprintf("%s!A2:J", tab))
	}

	resp, err := svc.Spreadsheets.Values.BatchGet(spreadsheetID).Ranges(ranges...).Do()
	if err != nil {
		return nil, err
	}

	index := make(map[string][]RewardActivity)
	for _, vr := range resp.ValueRanges {
		for _, row := range vr.Values {
			rollNorm := normalizeRollNo(safeGet(row, 3))
			if rollNorm == "" {
				continue
			}
			index[rollNorm] = append(index[rollNorm], RewardActivity{
				Date:         safeGet(row, 1),
				RewardPoints: safeGet(row, 7),
				ActivityType: safeGet(row, 8),
				ActivityName: safeGet(row, 9),
			})
		}
	}

	return index, nil
}

func (h *SheetHandler) getRewardsIndex() (map[string][]RewardActivity, error) {
	now := time.Now()

	rewardsCacheMu.RLock()
	if rewardsByRollCache != nil && now.Before(rewardsCacheExpiry) {
		cached := cloneRewardMap(rewardsByRollCache)
		rewardsCacheMu.RUnlock()
		return cached, nil
	}
	rewardsCacheMu.RUnlock()

	rewardsCacheBuildMu.Lock()
	defer rewardsCacheBuildMu.Unlock()

	rewardsCacheMu.RLock()
	if rewardsByRollCache != nil && time.Now().Before(rewardsCacheExpiry) {
		cached := cloneRewardMap(rewardsByRollCache)
		rewardsCacheMu.RUnlock()
		return cached, nil
	}
	rewardsCacheMu.RUnlock()

	index, err := h.buildRewardsIndex()
	if err != nil {
		rewardsCacheMu.RLock()
		stale := cloneRewardMap(rewardsByRollCache)
		rewardsCacheMu.RUnlock()
		if stale != nil {
			return stale, nil
		}
		return nil, err
	}

	rewardsCacheMu.Lock()
	rewardsByRollCache = index
	rewardsCacheExpiry = time.Now().Add(rewardsCacheTTL)
	rewardsCacheMu.Unlock()

	return cloneRewardMap(index), nil
}

func (h *SheetHandler) GetRewardsByRollNo(c *gin.Context) {
	rawQuery := strings.TrimSpace(c.Query("roll_no"))
	if rawQuery == "" {
		rawQuery = strings.TrimSpace(c.Query("rollno"))
	}
	if rawQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Query param 'roll_no' is required",
			"example": "/rewards?roll_no=7376251CS221",
		})
		return
	}

	queryNorm := normalizeRollNo(rawQuery)
	index, err := h.getRewardsIndex()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read reward sheet data",
		})
		return
	}

	matched := index[queryNorm]

	if len(matched) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No reward activities found for this roll number",
		})
		return
	}

	c.JSON(http.StatusOK, matched)
}
