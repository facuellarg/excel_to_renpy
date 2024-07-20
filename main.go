package main

import "os"

func main() {
	rows, err := ReadRenpyInfo("test.xlsx", "start")
	if err != nil {
		panic(err)
	}

	writer := NewWriter("./templates/excel_to_renpy.tmpl")
	text, err := writer.RenpyInfoToText(rows)
	if err != nil {
		panic(err)
	}
	print(text)
	//save text in renpy file
	if err := os.WriteFile("renpy.rpy", []byte(text), 0644); err != nil {
		panic(err)
	}

}
