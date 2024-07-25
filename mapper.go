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
	hideSplitChar    string
}

func NewDefaultMapper() *Mapper {
	return &Mapper{
		optionsSplitChar: "|",
		labelSplitChar:   ";",
		hideSplitChar:    "|",
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
	charactersBeingShown := make(map[string]string) //character being shown an the current expression

	for i, sheet := range sheets {
		rows := sheet.Rows
		if len(rows) == 0 {
			return nil, nil
		}

		renpyInfo.Labels = append(renpyInfo.Labels, models.Label{
			Label: sheet.Name,
			Scenes: []models.Scene{
				{Scene: "", Commands: []models.Command{}},
			},
		})

		currentScene := &renpyInfo.Labels[i].Scenes[0]

		for _, row := range rows {

			if row.Hide != "" {
				charactersToHide, _ := m.ParseHide(row.Hide)
				for _, character := range charactersToHide {
					currentScene.Commands = append(currentScene.Commands, models.Hide{
						Text: character,
					})
					delete(charactersBeingShown, character)
				}
			}

			if row.Kind == models.SceneKind {
				newScene := models.Scene{
					Scene:    row.Image,
					Commands: []models.Command{},
				}
				renpyInfo.Labels[i].Scenes = append(renpyInfo.Labels[i].Scenes, newScene)
				currentScene = &renpyInfo.Labels[i].Scenes[len(renpyInfo.Labels[i].Scenes)-1]
				charactersBeingShown = make(map[string]string)
				continue
			}

			if row.Kind == models.DialogueKind {
				if _, ok := charactersSet[row.Character]; !ok {
					charactersSet[row.Character] = struct{}{}
					renpyInfo.Characters = append(renpyInfo.Characters, row.Character)
				}

				if expression, ok := charactersBeingShown[row.Character]; !ok || expression != row.Expression {
					currentScene.Commands = append(currentScene.Commands, models.Show{
						Character:  row.Character,
						Position:   row.Position,
						Expression: row.Expression,
					})
					charactersBeingShown[row.Character] = row.Expression
				}

				currentScene.Commands = append(currentScene.Commands, models.Dialogue{
					Character: row.Character,
					Dialogue:  row.Text,
				})
				continue
			}

			if row.Kind == models.MenuKind {
				options, err := m.ParseOptions(row.Options)
				if err != nil {
					return nil, err
				}
				currentScene.Commands = append(currentScene.Commands, models.Menu{
					Options: options,
				})
			}

		}

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

func (m *Mapper) ParseHide(hide string) ([]string, error) {
	if hide == "" {
		return nil, nil
	}
	hideSplit := strings.Split(hide, m.hideSplitChar)
	return hideSplit, nil
}
