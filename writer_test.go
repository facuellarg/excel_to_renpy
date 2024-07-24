package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	tableTest := []struct {
		name         string
		sheetsInfo   []SheetInfo
		textExpected string
		errExpected  error
	}{
		{
			name: "Writes a Renpy file",
			sheetsInfo: []SheetInfo{
				{
					Name: "start",
					Rows: []RowInfo{
						{DialogueKind, "John", "Hello", "happy", "left", "", "", ""},
						{DialogueKind, "Tom", "How are you?", "happy", "left", "", "", ""},
						{SceneKind, "", "", "", "", "", "imageScene", ""},
						{DialogueKind, "John", "Hello in scene2", "happy", "left", "", "", ""},
					},
				},
			},
			textExpected: `define John = Character("John")
define Tom = Character("Tom")

label start:

  John "Hello"
  Tom "How are you?"

  scene imageScene
  John "Hello in scene2"`,
			errExpected: nil,
		},
		{
			name: "Writes a Renpy file with menu",
			sheetsInfo: []SheetInfo{
				{
					Name: "start",
					Rows: []RowInfo{
						{DialogueKind, "John", "Hello", "happy", "left", "", "", ""},
						{DialogueKind, "Tom", "How are you?", "happy", "left", "", "", ""},
						{SceneKind, "", "", "", "", "", "imageScene", ""},
						{DialogueKind, "John", "Hello in scene2", "happy", "left", "", "", ""},
						{MenuKind, "", "", "", "", "option1;otherLabel|option2|option3", "", ""},
					},
				},
			},
			textExpected: `define John = Character("John")
define Tom = Character("Tom")

label start:

  John "Hello"
  Tom "How are you?"

  scene imageScene
  John "Hello in scene2"
  menu:
    "option1":
      jump otherLabel
    "option2"
    "option3"`,
			errExpected: nil,
		},
		{
			name: "Writes a Renpy file with two sheets",
			sheetsInfo: []SheetInfo{
				{
					Name: "start",
					Rows: []RowInfo{
						{DialogueKind, "John", "Hello", "happy", "left", "", "", ""},
						{DialogueKind, "Tom", "How are you?", "happy", "left", "", "", ""},
						{SceneKind, "", "", "", "", "", "imageScene", ""},
						{DialogueKind, "John", "Hello in scene2", "happy", "left", "", "", ""},
						{MenuKind, "", "", "", "", "option1;otherLabel|option2|option3", "", ""},
					},
				},
				{
					Name: "anotherLabel",
					Rows: []RowInfo{
						{DialogueKind, "John", "Hello in another label", "happy", "left", "", "", ""},
					},
				},
			},
			textExpected: `define John = Character("John")
define Tom = Character("Tom")

label start:

  John "Hello"
  Tom "How are you?"

  scene imageScene
  John "Hello in scene2"
  menu:
    "option1":
      jump otherLabel
    "option2"
    "option3"

label anotherLabel:

  John "Hello in another label"`,
			errExpected: nil,
		},
	}
	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			renpyInfo := tt.sheetsInfo
			writer := NewWriter("./templates/excel_to_renpy.tmpl")
			text, err := writer.RenpyInfoToText(renpyInfo)
			assert.ErrorIs(t, err, tt.errExpected)
			assert.Equal(t, tt.textExpected, strings.TrimSpace(text))
		})
	}
}
