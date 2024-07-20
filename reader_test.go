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
				{"tester1", "dialog 1", "first scene"},
				{"tester2", "dialog 2", "first scene"},
			},
			errExpected: nil,
		},
	}
	for _, tt := range tableTest {

		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			renpyInfo, err := ReadRenpyInfo(path)

			assert.ErrorIs(t, err, tt.errExpected)
			// if err != nil {
			// 	t.Errorf("Expected to read Renpy file, got %v", err)
			// 	return
			// }
			assert.Equal(t, len(tt.renpyExpected), len(renpyInfo))
			for i, renpy := range renpyInfo {
				assert.Equal(t, tt.renpyExpected[i].Character, renpy.Character)
				assert.Equal(t, tt.renpyExpected[i].Dialogue, renpy.Dialogue)
				assert.Equal(t, tt.renpyExpected[i].Scene, renpy.Scene)
			}
		})
	}
}
