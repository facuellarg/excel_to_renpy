package main

import (
	"bytes"
	"log"
)

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

type Menu struct {
	Options []Options
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

func (o Menu) Build() string {
	var buffer bytes.Buffer
	if len(o.Options) == 0 {
		return ""
	}
	buffer.WriteString("menu:")
	for _, option := range o.Options {
		buffer.WriteString("\n")
		buffer.WriteString("    \"" + option.Text + "\"")
		if option.Label != "" {
			buffer.WriteString(":\n      jump " + option.Label)
		}
	}
	return buffer.String()
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
