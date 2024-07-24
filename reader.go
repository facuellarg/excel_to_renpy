package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Header string

const (
	KIND       Header = "kind"
	CHARACTER  Header = "character"
	TEXT       Header = "text"
	EXPRESSION Header = "expression"
	POSITION   Header = "position"
	OPTIONS    Header = "options"
	IMAGE      Header = "image"
	ANIMATION  Header = "animation"
)

var (
	HEADERS = map[Header]int{
		KIND:       0,
		CHARACTER:  1,
		TEXT:       2,
		EXPRESSION: 3,
		POSITION:   4,
		OPTIONS:    5,
		IMAGE:      6,
		ANIMATION:  7,
	}
)

func ReadExcel(path string) ([]SheetInfo, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the file
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	sheetNames := f.GetSheetList()
	sheetInfos := make([]SheetInfo, len(sheetNames))
	for i, sheet := range sheetNames {
		rows, err := ReadSheetInfo(f, sheet)
		if err != nil {
			return nil, err
		}
		sheetInfos[i] = SheetInfo{
			Name: sheet,
			Rows: rows,
		}
	}

	return sheetInfos, nil
}

func ReadSheetInfo(f *excelize.File, sheet string) ([]RowInfo, error) {

	rows, err := f.GetRows(sheet, excelize.Options{
		RawCellValue: true,
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found")
	}
	for j, header := range rows[0] {
		h := Header(header)
		if _, ok := HEADERS[h]; !ok {
			return nil, fmt.Errorf("header %s not found", header)
		}
		HEADERS[h] = j
	}

	renpyInfos := make([]RowInfo, len(rows)-1)
	for i, row := range rows[1:] {
		renpyInfos[i] = RowInfo{
			Kind:       StringToKind(GetValue(row, HEADERS[KIND])),
			Character:  GetValue(row, HEADERS[CHARACTER]),
			Text:       GetValue(row, HEADERS[TEXT]),
			Expression: GetValue(row, HEADERS[EXPRESSION]),
			Position:   GetValue(row, HEADERS[POSITION]),
			Options:    GetValue(row, HEADERS[OPTIONS]),
			Image:      GetValue(row, HEADERS[IMAGE]),
			Animation:  GetValue(row, HEADERS[ANIMATION]),
		}
	}

	return renpyInfos, nil
}

func GetValue[T any](row []T, index int) T {
	if index >= len(row) {
		var zero T
		return zero
	}
	return row[index]
}

func GetValueOrDefault[T any](row []T, index int, defaultValue T) T {
	if index >= len(row) {
		return defaultValue
	}
	return row[index]
}
