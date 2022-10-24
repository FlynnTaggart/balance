package utils

import (
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
	return strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-01"
}

func LastDayInMonth(year, month int) string {
	switch month {
	case 2:
		if year%400 == 0 || (year%100 != 0 && year%4 == 0) {
			return strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-29"
		}
		return strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-28"
	case 4:
		fallthrough
	case 6:
		fallthrough
	case 9:
		fallthrough
	case 11:
		return strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-30"
	default:
		return strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-31"
	}
}

func MoneyToInt(amount float32) int64 {
	return int64(math.Ceil(float64(amount * 100)))
}

func MoneyToFloat(amount int64) float32 {
	return float32(amount) * 0.01
}
