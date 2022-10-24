package utils

import (
	"fmt"
	"math"
	"strconv"
)

func GetReportFilePath(year int, month int) string {
	filePath := "./report/" + strconv.Itoa(year) + "/" + strconv.Itoa(month) + "/report.csv"
	return filePath
}

func GetReportFileDir(year int, month int) string {
	filePath := "./report/" + strconv.Itoa(year) + "/" + strconv.Itoa(month) + "/"
	return filePath
}

func FirstDayInMonth(year, month int) string {
	return strconv.Itoa(year) + "-" + fmt.Sprintf("%02d", month) + "-01"
}

func LastDayInMonth(year, month int) string {
	if month == 12 {
		return strconv.Itoa(year+1) + "-01-01"
	} else {
		return strconv.Itoa(year) + "-" + fmt.Sprintf("%02d", month+1) + "-01"
	}
}

func MoneyToInt(amount float32) int64 {
	return int64(math.Ceil(float64(amount * 100)))
}

func MoneyToFloat(amount int64) float32 {
	return float32(amount) * 0.01
}
