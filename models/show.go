package models

import "bytes"

type Show struct {
	Character  string
	Position   string
	Expression string
}

func (s Show) Build() string {
	var buffer bytes.Buffer
	buffer.WriteString("show ")
	buffer.WriteString(s.Character)
	if s.Expression != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s.Expression)
	}

	if s.Position != "" {
		buffer.WriteString(" at ")
		buffer.WriteString(s.Position)
	}

	return buffer.String()
}
