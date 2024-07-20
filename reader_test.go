package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRenpyInfo(t *testing.T) {

	tableTest := []struct {
		name          string
		path          string
		renpyExpected []RowInfo
		errExpected   error
	}{
		{
			name: "Reads a Renpy file",
			path: "test.xlsx",
			renpyExpected: []RowInfo{
				{DialogueKind, "John", "Hello", "happy", "left", "", "", ""},
				{DialogueKind, "John", "How are you?", "happy", "left", "", "", ""},
				{MenuKind, "", "", "", "", "option1|otherLabel;option2;option3", "", ""},
				{SceneKind, "", "", "", "", "", "imageScene", ""},
			},
			errExpected: nil,
		},
	}
	for _, tt := range tableTest {

		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			renpyInfo, err := ReadRenpyInfo(path, "start")

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
