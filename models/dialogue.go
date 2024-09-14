package models

type Dialogue struct {
	Character string
	Dialogue  string
}

func (d Dialogue) Build() string {
	return d.Character + " " + "\"" + d.Dialogue + "\""
}
