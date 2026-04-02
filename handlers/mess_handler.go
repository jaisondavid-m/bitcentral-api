package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/sheets/v4"
)

// Meal timing windows (IST — adjust as needed)
var mealTimings = map[string][2]string{
	"Breakfast": {"07:00", "08:30"},
	"Lunch":     {"12:20", "13:30"},
	"Dinner":    {"19:00", "20:30"},
}

// Hostel → Spreadsheet ID env var mapping
var hostelSheetEnv = map[string]string{
	"boys":  "BOYS_MESS_SPREADSHEET_ID",
	"girls": "GIRLS_MESS_SPREADSHEET_ID",
}

type MessHandler struct {
	sh *SheetHandler
}

func NewMessHandler(sh *SheetHandler) *MessHandler {
	return &MessHandler{sh: sh}
}


func currentMealType() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	hhmm := now.Format("15:04")

	order := []string{"Breakfast", "Lunch", "Dinner"}
	for _, meal := range order {
		window := mealTimings[meal]
		if hhmm >= window[0] && hhmm <= window[1] {
			return meal
		}
	}
	for _, meal := range order {
		if hhmm < mealTimings[meal][0] {
			return meal
		}
	}
	return "Breakfast"
}

func (h *MessHandler) fetchMessMenu(hostel, date string) (map[string][]string, string, error) {
	svc := h.sh.getSheetsService()
	if svc == nil {
		return nil, "", fmt.Errorf("sheets service not initialized")
	}

	envKey, ok := hostelSheetEnv[strings.ToLower(hostel)]
	if !ok {
		return nil, "", fmt.Errorf("unknown hostel '%s'; use 'boys' or 'girls'", hostel)
	}
	spreadsheetID := os.Getenv(envKey)
	if spreadsheetID == "" {
		return nil, "", fmt.Errorf("env var %s is not set", envKey)
	}

	resp, err := svc.Spreadsheets.Values.Get(spreadsheetID, "Sheet1!A1:D9999").Do()
	if err != nil {
		return nil, "", fmt.Errorf("sheets API error: %w", err)
	}

	menu := map[string][]string{
		"Breakfast": {},
		"Lunch":     {},
		"Dinner":    {},
	}
	dayName := ""

	for _, row := range resp.Values {
		if len(row) < 4 {
			continue
		}
		rowDate := strings.TrimSpace(fmt.Sprintf("%v", row[0]))
		if rowDate != date {
			continue
		}
		if dayName == "" {
			dayName = strings.TrimSpace(fmt.Sprintf("%v", row[1]))
		}
		mealType := strings.TrimSpace(fmt.Sprintf("%v", row[2]))
		item := strings.TrimSpace(fmt.Sprintf("%v", row[3]))
		if item == "" {
			continue
		}
		switch mealType {
		case "Breakfast", "Lunch", "Dinner":
			menu[mealType] = append(menu[mealType], item)
		}
	}

	return menu, dayName, nil
}

func (h *MessHandler) GetMess(c *gin.Context) {
	hostel := strings.ToLower(strings.TrimSpace(c.Query("hostel")))
	if hostel == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "'hostel' query param is required",
			"example": "/mess?hostel=boys&date=2026-04-02",
		})
		return
	}

	// Resolve date
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)

	dateStr := strings.TrimSpace(c.Query("date"))
	if dateStr == "" {
		dateStr = now.Format("2006-01-02")
	} else {
		if _, err := time.Parse("2006-01-02", dateStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid date format; use YYYY-MM-DD",
			})
			return
		}
	}

	menu, dayName, err := h.fetchMessMenu(hostel, dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activeMeal := currentMealType()
	window := mealTimings[activeMeal]

	// Build the current-meal section
	currentMealData := gin.H{
		"meal_type":  activeMeal,
		"start_time": window[0],
		"end_time":   window[1],
		"items":      menu[activeMeal],
	}

	// Full day menu
	fullMenu := gin.H{
		"breakfast": menu["Breakfast"],
		"lunch":     menu["Lunch"],
		"dinner":    menu["Dinner"],
	}

	hasData := len(menu["Breakfast"]) > 0 || len(menu["Lunch"]) > 0 || len(menu["Dinner"]) > 0

	c.JSON(http.StatusOK, gin.H{
		"hostel":       hostel,
		"date":         dateStr,
		"day":          dayName,
		"current_time": now.Format("15:04") + " IST",
		"current_meal": currentMealData,
		"full_menu":    fullMenu,
		"data_found":   hasData,
	})
}

// GetMealTimings handles: GET /mess/timings
// Returns the configured meal time windows.
func (h *MessHandler) GetMealTimings(c *gin.Context) {
	type window struct {
		Meal      string `json:"meal"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	timings := []window{
		{"Breakfast", mealTimings["Breakfast"][0], mealTimings["Breakfast"][1]},
		{"Lunch", mealTimings["Lunch"][0], mealTimings["Lunch"][1]},
		{"Dinner", mealTimings["Dinner"][0], mealTimings["Dinner"][1]},
	}
	c.JSON(http.StatusOK, gin.H{"timings": timings})
}

// Ensure sheets import is used (already used via SheetHandler).
var _ = sheets.SpreadsheetsReadonlyScope