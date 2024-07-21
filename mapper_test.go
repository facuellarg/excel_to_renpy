package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowsInfoToRenpyInfo(t *testing.T) {
	tt := []struct {
		name        string
		rowsInfo    []RowInfo
		renpyInfo   *RenpyInfo
		errExpected error
	}{
		{
			name: "Converts rows to RenpyInfo",
			rowsInfo: []RowInfo{
				{DialogueKind, "John", "Hello", "happy", "left", "", "", ""},
				{DialogueKind, "Tom", "How are you?", "happy", "left", "", "", ""},
				{MenuKind, "", "", "", "", "option1;otherLabel|option2|option3", "", ""},
				{SceneKind, "", "", "", "", "", "imageScene", ""},
				{DialogueKind, "John", "Hello in scene2", "happy", "left", "", "", ""},
			},
			renpyInfo: &RenpyInfo{
				Characters: []string{"John", "Tom"},
				Labels: []Label{
					{
						Label: "start", Scenes: []Scene{
							{"", []Command{
								Dialogue{"John", "Hello"},
								Dialogue{"Tom", "How are you?"},
								Menu{Options: []Options{
									{"option1", "otherLabel"},
									{"option2", ""},
									{"option3", ""},
								}},
							},
							},
							{"imageScene", []Command{Dialogue{"John", "Hello in scene2"}}},
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
		rowInfo     RowInfo
		expected    bool
		errExpected error
	}{
		{
			name:        "Dialogue is valid",
			rowInfo:     RowInfo{DialogueKind, "John", "Hello", "happy", "left", "", "", ""},
			expected:    true,
			errExpected: nil,
		},
		{
			name:        "Kind is not Dialogue",
			rowInfo:     RowInfo{MenuKind, "John", "Hello", "happy", "left", "", "", ""},
			expected:    false,
			errExpected: nil,
		},
		{
			name:        "Character is empty",
			rowInfo:     RowInfo{DialogueKind, "", "Hello", "happy", "left", "", "", ""},
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
		rowInfo     RowInfo
		expected    bool
		errExpected error
	}{
		{
			name:        "Scene is valid",
			rowInfo:     RowInfo{SceneKind, "", "", "", "", "", "imageScene", ""},
			expected:    true,
			errExpected: nil,
		},
		{
			name:        "Scene with empty image",
			rowInfo:     RowInfo{SceneKind, "", "", "", "", "", "", ""},
			expected:    false,
			errExpected: nil,
		},
		{
			name:        "Kind is not Scene",
			rowInfo:     RowInfo{MenuKind, "", "", "", "", "", "imageScene", ""},
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
		rowInfo     RowInfo
		expected    bool
		errExpected error
	}{
		{
			name:        "Menu is valid",
			rowInfo:     RowInfo{MenuKind, "", "", "", "", "option1|otherLabel;option2;option3", "", ""},
			expected:    true,
			errExpected: nil,
		},
		{
			name:        "Options is empty",
			rowInfo:     RowInfo{MenuKind, "", "", "", "", "", "", ""},
			expected:    false,
			errExpected: nil,
		},
		{
			name:        "Kind is not Menu",
			rowInfo:     RowInfo{SceneKind, "", "", "", "", "option1|otherLabel;option2;option3", "", ""},
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
		expected    []Options
		errExpected error
	}{
		{
			name:     "Parses options",
			options:  "option1;otherLabel|option2|option3",
			expected: []Options{{"option1", "otherLabel"}, {"option2", ""}, {"option3", ""}},
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
