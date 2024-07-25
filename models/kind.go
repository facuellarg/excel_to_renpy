package models

import "log"

type Kind int

//go:generate stringer -type=Kind
const (
	MenuKind Kind = iota
	SceneKind
	DialogueKind
)

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
