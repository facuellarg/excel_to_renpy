package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Header string

const (
	CHARACTER    Header = "character"
	EXPRESSION   Header = "expression"
	EMOTION      Header = "emotion"
	BACKGROUND   Header = "background"
	LOCATION     Header = "location"
	DIALOGUE     Header = "dialogue"
	SCENE        Header = "scene"
	DIALOGUETYPE Header = "dialoguetype"
)

var (
	HEADERS = map[Header]int{
		CHARACTER:    1,
		EXPRESSION:   2,
		EMOTION:      3,
		BACKGROUND:   4,
		LOCATION:     5,
		DIALOGUE:     6,
		SCENE:        7,
		DIALOGUETYPE: 8,
	}
)

func ReadRenpyInfo(path string) ([]RowInfo, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Close the file
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found")
	}
	for j, header := range rows[0] {
		h := Header(header)
		if _, ok := HEADERS[h]; !ok {
			return nil, fmt.Errorf("header %s not found", header)
		}
		HEADERS[h] = j
	}

	renpyInfos := make([]RowInfo, len(rows)-1)
	for i, row := range rows[1:] {
		renpyInfos[i] = RowInfo{
			Character: row[HEADERS[CHARACTER]],
			// Expression:   row[1],
			// Emotion:      row[2],
			// Background:   row[3],
			// Location:     row[4],
			Dialogue: row[HEADERS[DIALOGUE]],
			Scene:    row[HEADERS[SCENE]],
			// DialogueType: row[6],
		}
	}

	return renpyInfos, nil
}

func RowInfoToRenpy(rows []RowInfo) RenpyInfo {
	renpyInfo := RenpyInfo{}
	if len(rows) == 0 {
		return renpyInfo
	}

	renpyInfo.Characters = make([]string, 0, len(rows))
	charactersSet := make(map[string]struct{})

	renpyInfo.Scenes = make([]Scene, 0, len(rows))
	charactersSet[rows[0].Character] = struct{}{}
	renpyInfo.Characters = append(renpyInfo.Characters, rows[0].Character)

	renpyInfo.Scenes = append(renpyInfo.Scenes, Scene{
		Scene: rows[0].Scene,
		Dialogues: []Dialogue{
			{Character: rows[0].Character, Dialogue: rows[0].Dialogue},
		},
	})
	sceneIndex := 0
	lastScene := rows[0].Scene
	var d Dialogue

	for _, row := range rows[1:] {
		if _, ok := charactersSet[row.Character]; !ok {
			charactersSet[row.Character] = struct{}{}
			renpyInfo.Characters = append(renpyInfo.Characters, row.Character)
		}

		d = Dialogue{Character: row.Character, Dialogue: row.Dialogue}
		if row.Scene != lastScene {
			sceneIndex++
			lastScene = row.Scene
			renpyInfo.Scenes = append(renpyInfo.Scenes, Scene{
				Scene:     row.Scene,
				Dialogues: []Dialogue{d},
			})
		} else {
			renpyInfo.Scenes[sceneIndex].Dialogues = append(renpyInfo.Scenes[sceneIndex].Dialogues, d)
		}
	}

	return renpyInfo
}
