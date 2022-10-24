package databases

import (
	"balance/internal/utils"

	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

// CreateReport creates report file and returns relative path to it
func (p PgxDB) CreateReport(year, month int) (string, error) {
	ctx := context.Background()

	// log error
	var err error
	defer func() {
		if err != nil {
			p.Logger.Log(ctx, pgx.LogLevelError, fmt.Sprintf("db: create report: %v", err), nil)
		}
	}()

	type parsedRow struct {
		ServiceName string
		Amount      int64
	}

	// get all operation from given month
	from := utils.FirstDayInMonth(year, month)
	to := utils.LastDayInMonth(year, month)
	rows, _ := p.Query(ctx, "select * from operations where service_id is not null and done_at between $1 and $2 order by service_name",
		from, to)
	defer rows.Close()

	// parse rows to row struct
	var parsedRows []parsedRow
	for rows.Next() {
		var r parsedRow
		var tempUint uint64 // we need this because if we select certain columns pgx returns error: "number of field descriptions must equal number of destinations, got 1 and 2"
		var tempTime time.Time
		err := rows.Scan(&tempUint, &tempUint, &tempUint, &r.ServiceName, &r.Amount, &tempTime)
		if err != nil {
			return "", err
		}
		parsedRows = append(parsedRows, r)
	}
	if err = rows.Err(); err != nil {
		return "", err
	}

	// convert parsed rows to csv table rows
	csvRows := make(map[string]int64)
	for _, r := range parsedRows {
		csvRows[r.ServiceName] += r.Amount
	}

	// open out file
	// if file and dir for it are not created - create them
	filePath := utils.GetReportFilePath(year, month)
	_, err = os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(utils.GetReportFileDir(year, month), os.ModePerm)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		err = file.Close()
	}()

	w := csv.NewWriter(file)
	defer w.Flush()

	// write header and rows of csv table
	err = w.Write([]string{"service_name", "month_amount"})
	if err != nil {
		return "", err
	}
	for name, amount := range csvRows {
		if err = w.Write([]string{name, fmt.Sprintf("%.2f", float32(amount)*-0.01)}); err != nil {
			return "", err
		}
	}

	return filePath[1:], nil
}
