package models

type Command interface {
	Build() string
}
