package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	tableTest := []struct {
		name         string
		renpyInfo    []RowInfo
		textExpected string
		errExpected  error
	}{
		{
			name: "Writes a Renpy file",
			renpyInfo: []RowInfo{
				{"tester1", "dialog 1", "first scene"},
				{"tester2", "dialog 2", "first scene"},
				{"tester1", "dialog 3", "second scene"},
			},
			textExpected: `define tester1 = Character("tester1")
define tester2 = Character("tester2")

label start:
	scene first scene
	tester1 "dialog 1"
	tester2 "dialog 2"
	scene second scene
	tester1 "dialog 3"
`,
			errExpected: nil,
		},
	}
	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			renpyInfo := tt.renpyInfo
			writer := NewWriter("./templates/excel_to_renpy.tmpl")
			text, err := writer.RenpyInfoToText(renpyInfo)
			print(text)
			assert.ErrorIs(t, err, tt.errExpected)
			assert.Equal(t, tt.textExpected, text)
		})
	}
}
