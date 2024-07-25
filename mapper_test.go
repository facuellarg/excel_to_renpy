package main

import (
	"renpy-transformer/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowsInfoToRenpyInfo(t *testing.T) {
	tt := []struct {
		name        string
		rowsInfo    []models.SheetInfo
		renpyInfo   *models.RenpyInfo
		errExpected error
	}{
		{
			name: "Converts rows to RenpyInfo",
			rowsInfo: []models.SheetInfo{
				{
					Name: "start",
					Rows: []models.RowInfo{
						{Kind: models.DialogueKind, Character: "John", Text: "Hello", Expression: "happy", Position: "left"},
						{Kind: models.DialogueKind, Character: "Tom", Text: "How are you?", Expression: "happy", Position: "left"},
						{Kind: models.DialogueKind, Character: "Tom", Text: "I am angry", Expression: "angry", Position: "left"},
						{Kind: models.DialogueKind, Character: "John", Text: "I still happy", Expression: "happy", Position: "left"},
						{Kind: models.MenuKind, Options: "option1;otherLabel|option2|option3", Hide: "Tom"},
						{Kind: models.SceneKind, Image: "imageScene"},
						{Kind: models.DialogueKind, Character: "John", Text: "Hello in scene2", Expression: "happy", Position: "left"},
					},
				},
				{
					Name: "otherLabel",
					Rows: []models.RowInfo{
						{Kind: models.DialogueKind, Character: "John", Text: "Hello in another label", Expression: "happy", Position: "left"},
					},
				},
			},
			renpyInfo: &models.RenpyInfo{
				Characters: []string{"John", "Tom"},
				Labels: []models.Label{
					{
						Label: "start", Scenes: []models.Scene{
							{Scene: "", Commands: []models.Command{
								models.Show{Character: "John", Expression: "happy", Position: "left"},
								models.Dialogue{Character: "John", Dialogue: "Hello"},
								models.Show{Character: "Tom", Expression: "happy", Position: "left"},
								models.Dialogue{Character: "Tom", Dialogue: "How are you?"},
								models.Show{Character: "Tom", Expression: "angry", Position: "left"},
								models.Dialogue{Character: "Tom", Dialogue: "I am angry"},
								models.Dialogue{Character: "John", Dialogue: "I still happy"},
								models.Hide{Text: "Tom"},
								models.Menu{Options: []models.Options{
									{Text: "option1", Label: "otherLabel"},
									{Text: "option2", Label: ""},
									{Text: "option3", Label: ""},
								}},
							},
							},
							{Scene: "imageScene", Commands: []models.Command{
								models.Show{Character: "John", Expression: "happy", Position: "left"},
								models.Dialogue{Character: "John", Dialogue: "Hello in scene2"},
							}},
						},
					},
					{
						Label: "otherLabel", Scenes: []models.Scene{
							{
								Scene: "", Commands: []models.Command{
									models.Dialogue{Character: "John", Dialogue: "Hello in another label"},
								},
							},
						},
					},
				},
			},
			errExpected: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mapper := NewDefaultMapper()
			renpyInfo, err := mapper.RowsInfoToRenpyInfo(tc.rowsInfo)
			if tc.errExpected == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, err, tc.errExpected)
			}
			if renpyInfo == nil {
				if tc.renpyInfo == nil {
					return
				}
				t.Fatalf("Expected RenpyInfo, got nil")
			}
			assert.Equal(t, tc.renpyInfo.Characters, renpyInfo.Characters)
			assert.Equal(t, len(tc.renpyInfo.Labels), len(renpyInfo.Labels))
			for i, label := range renpyInfo.Labels {
				assert.Equal(t, tc.renpyInfo.Labels[i].Label, label.Label)
				sceneExpected := tc.renpyInfo.Labels[i].Scenes
				assert.Equal(t, len(sceneExpected), len(label.Scenes))
				for i, scene := range label.Scenes {
					assert.Equal(t, sceneExpected[i].Scene, scene.Scene)
					assert.Equal(t, len(sceneExpected[i].Commands), len(scene.Commands))
					for j, command := range scene.Commands {
						assert.Equal(t, sceneExpected[i].Commands[j].Build(), command.Build())
					}
				}
			}
		})
	}

}

func TestValidateDialogue(t *testing.T) {

	tt := []struct {
		name        string
		rowInfo     models.RowInfo
		expected    bool
		errExpected error
	}{
		{
			name:        "Dialogue is valid",
			rowInfo:     models.RowInfo{Kind: models.DialogueKind, Character: "John", Text: "Hello", Expression: "happy", Position: "left", Options: "", Image: "", Animation: ""},
			expected:    true,
			errExpected: nil,
		},
		{
			name:        "Kind is not Dialogue",
			rowInfo:     models.RowInfo{Kind: models.MenuKind, Character: "John", Text: "Hello", Expression: "happy", Position: "left", Options: "", Image: "", Animation: ""},
			expected:    false,
			errExpected: nil,
		},
		{
			name:        "Character is empty",
			rowInfo:     models.RowInfo{Kind: models.DialogueKind, Character: "", Text: "Hello", Expression: "happy", Position: "left", Options: "", Image: "", Animation: ""},
			expected:    false,
			errExpected: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mapper := Mapper{}
			result, err := mapper.ValidateDialog(tc.rowInfo)
			if tc.errExpected == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorAs(t, err, tc.errExpected)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateScene(t *testing.T) {
	tt := []struct {
		name        string
		rowInfo     models.RowInfo
		expected    bool
		errExpected error
	}{
		{
			name:        "Scene is valid",
			rowInfo:     models.RowInfo{Kind: models.SceneKind, Character: "", Text: "", Expression: "", Position: "", Options: "", Image: "imageScene", Animation: ""},
			expected:    true,
			errExpected: nil,
		},
		{
			name:        "Scene with empty image",
			rowInfo:     models.RowInfo{Kind: models.SceneKind, Character: "", Text: "", Expression: "", Position: "", Options: "", Image: "", Animation: ""},
			expected:    false,
			errExpected: nil,
		},
		{
			name:        "Kind is not Scene",
			rowInfo:     models.RowInfo{Kind: models.MenuKind, Character: "", Text: "", Expression: "", Position: "", Options: "", Image: "imageScene", Animation: ""},
			expected:    false,
			errExpected: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mapper := Mapper{}
			result, err := mapper.ValidateScene(tc.rowInfo)
			if tc.errExpected == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorAs(t, err, tc.errExpected)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateMenu(t *testing.T) {
	tt := []struct {
		name        string
		rowInfo     models.RowInfo
		expected    bool
		errExpected error
	}{
		{
			name:        "Menu is valid",
			rowInfo:     models.RowInfo{Kind: models.MenuKind, Character: "", Text: "", Expression: "", Position: "", Options: "option1|otherLabel;option2;option3", Image: "", Animation: ""},
			expected:    true,
			errExpected: nil,
		},
		{
			name:        "Options is empty",
			rowInfo:     models.RowInfo{Kind: models.MenuKind, Character: "", Text: "", Expression: "", Position: "", Options: "", Image: "", Animation: ""},
			expected:    false,
			errExpected: nil,
		},
		{
			name:        "Kind is not Menu",
			rowInfo:     models.RowInfo{Kind: models.SceneKind, Character: "", Text: "", Expression: "", Position: "", Options: "option1|otherLabel;option2;option3", Image: "", Animation: ""},
			expected:    false,
			errExpected: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mapper := Mapper{}
			result, err := mapper.ValidateMenu(tc.rowInfo)
			if tc.errExpected == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorAs(t, err, tc.errExpected)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseOptions(t *testing.T) {
	tt := []struct {
		name        string
		options     string
		expected    []models.Options
		errExpected error
	}{
		{
			name:     "Parses options",
			options:  "option1;otherLabel|option2|option3",
			expected: []models.Options{{Text: "option1", Label: "otherLabel"}, {Text: "option2", Label: ""}, {Text: "option3", Label: ""}},
		},
		{
			name:        "Invalid options",
			options:     "option;label;label|option2",
			expected:    nil,
			errExpected: ErrInvalidOptions,
		},
		{
			name:        "Options is empty",
			options:     "",
			expected:    nil,
			errExpected: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mapper := NewDefaultMapper()
			result, err := mapper.ParseOptions(tc.options)
			if tc.errExpected == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, err, tc.errExpected)
			}
			assert.Equal(t, len(tc.expected), len(result))
			for i, option := range result {
				assert.Equal(t, tc.expected[i], option)
			}
		})
	}
}
