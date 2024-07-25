package main

import (
	"renpy-transformer/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func TestReadExcel(t *testing.T) {

	tableTest := []struct {
		name              string
		path              string
		sheetInfoExpected []models.SheetInfo
		errExpected       error
	}{
		{
			name: "Reads a Renpy file with two sheets",
			path: "test.xlsx",
			sheetInfoExpected: []models.SheetInfo{
				{
					Name: "start",
					Rows: []models.RowInfo{
						{Kind: models.DialogueKind, Character: "John", Text: "Hello", Expression: "happy", Position: "left"},
						{Kind: models.DialogueKind, Character: "Tom", Text: "How are you?", Expression: "happy", Position: "left"},
						{Kind: models.MenuKind, Options: "option1;otherLabel|option2|option3"},
						{Kind: models.SceneKind, Image: "imageScene", Hide: "John"},
						{Kind: models.DialogueKind, Character: "Tom", Text: "Hello after scene", Expression: "happy", Position: "left"},
						{Kind: models.DialogueKind, Character: "Tom", Text: "I am angry now", Expression: "angry", Position: "left"},
					},
				},
				{
					Name: "otherLabel",
					Rows: []models.RowInfo{
						{Kind: models.DialogueKind, Character: "Tom", Text: "Hello from another label", Expression: "angry", Position: "left", Options: "", Image: "", Animation: ""},
					},
				},
			},
			errExpected: nil,
		},
	}
	for _, tt := range tableTest {

		t.Run(tt.name, func(t *testing.T) {
			renpyInfo, err := ReadExcel(tt.path)

			assert.ErrorIs(t, err, tt.errExpected)
			// if err != nil {
			// 	t.Errorf("Expected to read Renpy file, got %v", err)
			// 	return
			// }
			assert.Equal(t, len(tt.sheetInfoExpected), len(renpyInfo))
			for i, renpy := range renpyInfo {
				assert.Equal(t, tt.sheetInfoExpected[i].Name, renpy.Name)
				for j, row := range renpy.Rows {
					assert.Equal(t, tt.sheetInfoExpected[i].Rows[j].Kind, row.Kind)
					assert.Equal(t, tt.sheetInfoExpected[i].Rows[j].Character, row.Character)
					assert.Equal(t, tt.sheetInfoExpected[i].Rows[j].Text, row.Text)
					assert.Equal(t, tt.sheetInfoExpected[i].Rows[j].Expression, row.Expression)
					assert.Equal(t, tt.sheetInfoExpected[i].Rows[j].Position, row.Position)
					assert.Equal(t, tt.sheetInfoExpected[i].Rows[j].Options, row.Options)
					assert.Equal(t, tt.sheetInfoExpected[i].Rows[j].Image, row.Image)
					assert.Equal(t, tt.sheetInfoExpected[i].Rows[j].Animation, row.Animation)
					assert.Equal(t, tt.sheetInfoExpected[i].Rows[j].Hide, row.Hide)
				}
			}
		})
	}
}

func TestReadSheetInfo(t *testing.T) {

	tableTest := []struct {
		name          string
		path          string
		renpyExpected []models.RowInfo
		errExpected   error
	}{
		{
			name: "Reads a Renpy file",
			path: "test.xlsx",
			renpyExpected: []models.RowInfo{
				{Kind: models.DialogueKind, Character: "John", Text: "Hello", Expression: "happy", Position: "left"},
				{Kind: models.DialogueKind, Character: "Tom", Text: "How are you?", Expression: "happy", Position: "left"},
				{Kind: models.MenuKind, Options: "option1;otherLabel|option2|option3"},
				{Kind: models.SceneKind, Image: "imageScene", Hide: "John"},
				{Kind: models.DialogueKind, Character: "Tom", Text: "Hello after scene", Expression: "happy", Position: "left"},
				{Kind: models.DialogueKind, Character: "Tom", Text: "I am angry now", Expression: "angry", Position: "left"},
			},
			errExpected: nil,
		},
	}
	for _, tt := range tableTest {

		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			f, err := excelize.OpenFile(path)
			if err != nil {
				t.Errorf("Expected to open file, got %v", err)
				return
			}
			renpyInfo, err := ReadSheetInfo(f, "start")

			assert.ErrorIs(t, err, tt.errExpected)
			// if err != nil {
			// 	t.Errorf("Expected to read Renpy file, got %v", err)
			// 	return
			// }
			assert.Equal(t, len(tt.renpyExpected), len(renpyInfo))
			for i, renpy := range renpyInfo {
				assert.Equal(t, tt.renpyExpected[i].Kind, renpy.Kind)
				assert.Equal(t, tt.renpyExpected[i].Character, renpy.Character)
				assert.Equal(t, tt.renpyExpected[i].Text, renpy.Text)
				assert.Equal(t, tt.renpyExpected[i].Expression, renpy.Expression)
				assert.Equal(t, tt.renpyExpected[i].Position, renpy.Position)
				assert.Equal(t, tt.renpyExpected[i].Options, renpy.Options)
				assert.Equal(t, tt.renpyExpected[i].Image, renpy.Image)
				assert.Equal(t, tt.renpyExpected[i].Animation, renpy.Animation)
			}
		})
	}
}
