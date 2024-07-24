package main

import (
	"errors"
	"renpy-transformer/models"
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
		optionsSplitChar: "|",
		labelSplitChar:   ";",
	}
}

func (m *Mapper) RowsInfoToRenpyInfo(sheets []models.SheetInfo) (*models.RenpyInfo, error) {
	if len(sheets) == 0 {
		return nil, nil
	}

	renpyInfo := models.RenpyInfo{}
	renpyInfo.Characters = make([]string, 0)
	charactersSet := make(map[string]struct{})
	renpyInfo.Labels = make([]models.Label, 0)

	for i, sheet := range sheets {
		rows := sheet.Rows
		if len(rows) == 0 {
			return nil, nil
		}

		renpyInfo.Labels = append(renpyInfo.Labels, models.Label{
			Label: sheet.Name,
		})

		currentScene := renpyInfo.Labels[i].Scenes
		for _, row := range rows {
			if row.Kind == models.SceneKind {
				currentScene = append(currentScene, models.Scene{
					Scene:    row.Image,
					Commands: []models.Command{},
				})
			}

			if row.Kind == models.DialogueKind {
				if len(currentScene) == 0 {
					// return nil, errors.New("No scene found")
					currentScene = append(currentScene, models.Scene{})
				}
				if _, ok := charactersSet[row.Character]; !ok {
					charactersSet[row.Character] = struct{}{}
					renpyInfo.Characters = append(renpyInfo.Characters, row.Character)
				}
				currentScene[len(currentScene)-1].Commands = append(currentScene[len(currentScene)-1].Commands, models.Dialogue{
					Character: row.Character,
					Dialogue:  row.Text,
				})
			}
			if row.Kind == models.MenuKind {
				options, err := m.ParseOptions(row.Options)
				if err != nil {
					return nil, err
				}
				currentScene[len(currentScene)-1].Commands = append(currentScene[len(currentScene)-1].Commands, models.Menu{
					Options: options,
				})
			}
		}
		renpyInfo.Labels[i].Scenes = currentScene
	}
	return &renpyInfo, nil

}

func (m *Mapper) ValidateDialog(row models.RowInfo) (bool, error) {
	if row.Kind != models.DialogueKind {
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

func (m *Mapper) ValidateScene(row models.RowInfo) (bool, error) {
	if row.Kind != models.SceneKind {
		return false, nil
	}

	if row.Image == "" {
		return false, nil
	}

	return true, nil
}

func (m *Mapper) ValidateMenu(row models.RowInfo) (bool, error) {
	if row.Kind != models.MenuKind {
		return false, nil
	}

	if row.Options == "" {
		return false, nil
	}
	return true, nil

}

func (m *Mapper) ParseOptions(options string) ([]models.Options, error) {

	if options == "" {
		return nil, nil
	}

	optionsSplit := strings.Split(options, m.optionsSplitChar)
	optionsList := make([]models.Options, 0, len(optionsSplit))

	for _, option := range optionsSplit {
		optionSplit := strings.Split(option, m.labelSplitChar)
		if len(optionSplit) > 2 {
			return nil, ErrInvalidOptions
		}

		op := models.Options{
			Text: optionSplit[0],
		}
		if len(optionSplit) == 2 {
			op.Label = optionSplit[1]
		}
		optionsList = append(optionsList, op)

	}

	return optionsList, nil
}
