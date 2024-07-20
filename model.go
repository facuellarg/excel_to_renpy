package main

type Kind int

//go:generate stringer -type=Kind
const (
	CharacterKind Kind = iota
	SceneKind
	DialogueKind
)

type RowInfo struct {
	Character string
	Dialogue  string
	Scene     string
}

type RenpyInfo struct {
	Characters []string
	Scenes     []Scene
}

type Scene struct {
	// Background string
	// Location   string
	Scene     string
	Dialogues []Dialogue
}

type Dialogue struct {
	Character string
	Dialogue  string
}
