package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	tableTest := []struct {
		name         string
		rowsInfo     []RowInfo
		textExpected string
		errExpected  error
	}{
		{
			name: "Writes a Renpy file",
			rowsInfo: []RowInfo{
				{DialogueKind, "John", "Hello", "happy", "left", "", "", ""},
				{DialogueKind, "Tom", "How are you?", "happy", "left", "", "", ""},
				{MenuKind, "", "", "", "", "option1|otherLabel;option2;option3", "", ""},
				{SceneKind, "", "", "", "", "", "imageScene", ""},
				{DialogueKind, "John", "Hello in scene2", "happy", "left", "", "", ""},
			},
			textExpected: `define John = Character("John")
define Tom = Character("Tom")

label start:

  John "Hello"
  Tom "How are you?"

  scene imageScene
  John "Hello in scene2"
`,
			errExpected: nil,
		},
	}
	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			renpyInfo := tt.rowsInfo
			writer := NewWriter("./templates/excel_to_renpy.tmpl")
			text, err := writer.RenpyInfoToText(renpyInfo)
			print(text)
			assert.ErrorIs(t, err, tt.errExpected)
			assert.Equal(t, tt.textExpected, text)
		})
	}
}
