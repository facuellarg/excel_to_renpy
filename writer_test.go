package main

import (
	"renpy-transformer/models"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	tableTest := []struct {
		name         string
		sheetsInfo   []models.SheetInfo
		textExpected string
		errExpected  error
	}{
		{
			name: "Writes a Renpy file",
			sheetsInfo: []models.SheetInfo{
				{
					Name: "start",
					Rows: []models.RowInfo{
						{Kind: models.DialogueKind, Character: "John", Text: "Hello", Expression: "happy", Position: "left"},
						{Kind: models.DialogueKind, Character: "Tom", Text: "How are you?", Expression: "happy", Position: "left"},
						{Kind: models.MenuKind, Options: "option1;otherLabel|option2|option3"},
						{Kind: models.SceneKind, Image: "imageScene", Hide: "Tom"},
						{Kind: models.DialogueKind, Character: "John", Text: "Hello in scene2", Expression: "happy", Position: "right"},
						{Kind: models.DialogueKind, Character: "John", Text: "I am angry", Expression: "angry", Position: "right"},
					},
				},
				{
					Name: "otherLabel",
					Rows: []models.RowInfo{
						{Kind: models.DialogueKind, Character: "John", Text: "Hello in another label", Expression: "angry", Position: "left"},
						{Kind: models.DialogueKind, Character: "Tom", Text: "Hello in another label", Expression: "happy", Position: "left"},
					},
				},
			},
			textExpected: `define John = Character("John")
define Tom = Character("Tom")

label start:

  show John happy at left
  John "Hello"
  show Tom happy at left
  Tom "How are you?"
  menu:
    "option1":
      jump otherLabel
    "option2"
    "option3"
  hide Tom

  scene imageScene
  show John happy at right
  John "Hello in scene2"
  show John angry at right
  John "I am angry"

label otherLabel:

  John "Hello in another label"
  show Tom happy at left
  Tom "Hello in another label"`,
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
