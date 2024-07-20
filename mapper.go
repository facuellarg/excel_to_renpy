package main

import (
	"errors"
	"strings"
)

var (
	ErrInvalidOptions = errors.New("Invalid options")
)

type Mapper struct {
	optionsSplitChar string
	labelSplitChar   string
}

func NewDefaultMapper() *Mapper {
	return &Mapper{
		optionsSplitChar: ";",
		labelSplitChar:   "|",
	}
}

func (m *Mapper) RowsInfoToRenpyInfo(rows []RowInfo) (*RenpyInfo, error) {
	if len(rows) == 0 {
		return nil, nil
	}
	renpyInfo := RenpyInfo{}

	renpyInfo.Characters = make([]string, 0, len(rows))
	charactersSet := make(map[string]struct{})
	renpyInfo.Labels = make([]Label, 0)
	renpyInfo.Labels = append(renpyInfo.Labels, Label{
		Label: "start",
	})
	for _, row := range rows {
		if row.Kind == SceneKind {
			renpyInfo.Labels[0].Scenes = append(renpyInfo.Labels[0].Scenes, Scene{
				Scene:     row.Image,
				Dialogues: []Dialogue{},
			})
		}

		if row.Kind == DialogueKind {
			if len(renpyInfo.Labels[0].Scenes) == 0 {
				// return nil, errors.New("No scene found")
				renpyInfo.Labels[0].Scenes = append(renpyInfo.Labels[0].Scenes, Scene{})
			}
			if _, ok := charactersSet[row.Character]; !ok {
				charactersSet[row.Character] = struct{}{}
				renpyInfo.Characters = append(renpyInfo.Characters, row.Character)
			}
			renpyInfo.Labels[0].Scenes[len(renpyInfo.Labels[0].Scenes)-1].Dialogues = append(renpyInfo.Labels[0].Scenes[len(renpyInfo.Labels[0].Scenes)-1].Dialogues, Dialogue{
				Character: row.Character,
				Dialogue:  row.Text,
			})
		}
	}
	return &renpyInfo, nil

}

func (m *Mapper) ValidateDialog(row RowInfo) (bool, error) {
	if row.Kind != DialogueKind {
		return false, nil
	}

	if row.Character == "" {
		return false, nil
	}

	if row.Text == "" {
		return false, nil
	}

	return true, nil
}

func (m *Mapper) ValidateScene(row RowInfo) (bool, error) {
	if row.Kind != SceneKind {
		return false, nil
	}

	if row.Image == "" {
		return false, nil
	}

	return true, nil
}

func (m *Mapper) ValidateMenu(row RowInfo) (bool, error) {
	if row.Kind != MenuKind {
		return false, nil
	}

	if row.Options == "" {
		return false, nil
	}
	return true, nil

}

func (m *Mapper) ParseOptions(options string) ([]Options, error) {

	if options == "" {
		return nil, nil
	}

	optionsSplit := strings.Split(options, m.optionsSplitChar)
	optionsList := make([]Options, 0, len(optionsSplit))

	for _, option := range optionsSplit {
		optionSplit := strings.Split(option, m.labelSplitChar)
		if len(optionSplit) > 2 {
			return nil, ErrInvalidOptions
		}

		op := Options{
			Text: optionSplit[0],
		}
		if len(optionSplit) == 2 {
			op.Label = optionSplit[1]
		}
		optionsList = append(optionsList, op)

	}

	return optionsList, nil
}
