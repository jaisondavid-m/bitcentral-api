package data

import "server/models"

var Holidays = []models.Holiday{
	{
		FromDate: "2026-04-14",
		ToDate:   "2026-04-14",
		Day:      "Tuesday",
		Name:     "Tamil New Year",
	},
	{
		FromDate: "2026-04-24 (AN)",
		ToDate:   "2026-04-28",
		Name:     "GP",
	},
	{
		FromDate: "2026-05-01",
		ToDate:   "2026-05-01",
		Day:      "Friday",
		Name:     "May Day",
	},
}