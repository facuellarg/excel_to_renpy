package main

import "log"

type Kind int

//go:generate stringer -type=Kind
const (
	MenuKind Kind = iota
	SceneKind
	DialogueKind
)

type RowInfo struct {
	Kind       Kind
	Character  string
	Text       string
	Expression string
	Position   string
	Options    string
	Image      string
	Animation  string
}

type RenpyInfo struct {
	Characters []string
	Labels     []Label
}

type Scene struct {
	Scene    string
	Commands []Command
}

type Label struct {
	Label  string
	Scenes []Scene
}

type Dialogue struct {
	Character string
	Dialogue  string
}

type Options struct {
	Text  string
	Label string
}

type Command interface {
	Build() string
}

func (d Dialogue) Build() string {
	return d.Character + " " + "\"" + d.Dialogue + "\""
}

func StringToKind(kind string) Kind {
	switch kind {
	case "menu":
		return MenuKind
	case "scene":
		return SceneKind
	case "dialogue":
		return DialogueKind
	default:
		log.Printf("Unknown kind %s", kind)
		return -1
	}
}
