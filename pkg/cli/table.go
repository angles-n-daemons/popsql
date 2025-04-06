package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/angles-n-daemons/popsql/pkg/db/execution"
)

var Reset = "\033[0m"
var Bold = "\033[1m"
var Green = "\033[32m"

func TableRender(result *execution.Result) string {
	rows := rowsToStrings(result.Rows)
	lengths := make([]int, len(result.Columns))

	for _, row := range append(rows, result.Columns) {
		for i, val := range row {
			if len(val) > lengths[i] {
				lengths[i] = len(val)
			}
		}
	}

	output := splitter(lengths)
	output += fHeader(result.Columns, lengths)
	output += splitter(lengths)
	output += fRows(rows, lengths)
	output += splitter(lengths)
	output += duration(result.Duration)

	return output
}

func rowsToStrings(rows []execution.Row) [][]string {
	strs := make([][]string, len(rows))
	for i, row := range rows {
		strs[i] = make([]string, len(row))
		for j, v := range row {
			strs[i][j] = fmt.Sprintf("%v", v)
		}
	}
	return strs
}

func fHeader(columns []string, lengths []int) string {
	return styledRow(Bold+Green, Reset, columns, lengths)
}

func fRows(rows [][]string, lengths []int) string {
	s := ""
	for _, row := range rows {
		s += fRow(row, lengths)
	}
	return s
}

func fRow(row []string, lengths []int) string {
	return styledRow("", "", row, lengths)
}

func styledRow(style, styleEnd string, row []string, lengths []int) string {
	s := delimiter(true, false)
	fValues := make([]string, len(row))
	for i, val := range row {
		pad := strings.Repeat(" ", lengths[i]-len(val))
		fValues[i] = style + val + pad + styleEnd
	}
	s += strings.Join(fValues, delimiter(false, false))
	return s + delimiter(false, true) + "\n"
}

func splitter(lengths []int) string {
	s := "+"
	for _, l := range lengths {
		s += strings.Repeat("-", l+2)
	}
	return s + "+" + "\n"
}

func delimiter(start, end bool) string {
	d := "|"
	if !start {
		d = " " + d
	}
	if !end {
		d = d + " "
	}
	return d
}

func duration(d time.Duration) string {
	s := d / time.Second
	ms := (d / time.Millisecond) % 1000
	us := (d / time.Microsecond) % 1000
	return fmt.Sprintf("%ds %dms %dus\n", s, ms, us) + "\n"
}
