package main

import "os"

func main() {
	rows, err := ReadExcel("test 2.xlsx")
	if err != nil {
		panic(err)
	}

	writer := NewWriter("./templates/excel_to_renpy.tmpl")
	text, err := writer.RenpyInfoToText(rows)
	if err != nil {
		panic(err)
	}
	//save text in renpy file
	if err := os.WriteFile("/Users/freddycuellar/Documents/renpy/testing/game/script.rpy", []byte(text), 0644); err != nil {
		panic(err)
	}
}
