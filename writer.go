package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type Writer struct {
	excelToRenpy *template.Template
}

func NewWriter(path string) *Writer {
	w := Writer{}
	w.excelToRenpy = template.Must(template.ParseFiles(path))
	return &w
}

func (w *Writer) RenpyInfoToText(rows []RowInfo) (string, error) {
	renpyInfo := RowInfoToRenpy(rows)
	fmt.Printf("%+v\n", renpyInfo)
	var buffer bytes.Buffer
	err := w.excelToRenpy.Execute(&buffer, renpyInfo)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
