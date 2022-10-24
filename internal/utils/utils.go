package utils

import (
	"fmt"
	"math"
	"strconv"
)

// GetReportFilePath returns filepath by given year and month
func GetReportFilePath(year int, month int) string {
	filePath := "./report/" + strconv.Itoa(year) + "/" + strconv.Itoa(month) + "/report.csv"
	return filePath
}

// GetReportFileDir returns dir path by given year and month
func GetReportFileDir(year int, month int) string {
	filePath := "./report/" + strconv.Itoa(year) + "/" + strconv.Itoa(month) + "/"
	return filePath
}

// FirstDayInMonth returns date string with first day in given month
func FirstDayInMonth(year, month int) string {
	return strconv.Itoa(year) + "-" + fmt.Sprintf("%02d", month) + "-01"
}

// LastDayInMonth returns date string with last day in given month (actually first day of the next month)
func LastDayInMonth(year, month int) string {
	if month == 12 {
		return strconv.Itoa(year+1) + "-01-01"
	} else {
		return strconv.Itoa(year) + "-" + fmt.Sprintf("%02d", month+1) + "-01"
	}
}

// MoneyToInt converts money amount from float to number of cents
func MoneyToInt(amount float32) int64 {
	return int64(math.Ceil(float64(amount * 100)))
}

// MoneyToFloat converts money amount from float to number of cents
func MoneyToFloat(amount int64) float32 {
	return float32(amount) * 0.01
}
