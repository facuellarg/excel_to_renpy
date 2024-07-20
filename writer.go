package main

import (
	"bytes"
	"text/template"
)

type Writer struct {
	excelToRenpy *template.Template
	mapper       *Mapper
}

func NewWriter(path string) *Writer {
	w := Writer{}
	w.excelToRenpy = template.Must(template.ParseFiles(path))
	w.mapper = NewDefaultMapper()
	return &w
}

func (w *Writer) RenpyInfoToText(rows []RowInfo) (string, error) {
	renpyInfo, err := w.mapper.RowsInfoToRenpyInfo(rows)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = w.excelToRenpy.Execute(&buffer, renpyInfo)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
