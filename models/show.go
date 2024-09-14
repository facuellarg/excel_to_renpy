package models

import (
	"bytes"
	"fmt"
	"strings"
)

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
		if strings.Contains(s.Expression, " ") {
			fmt.Printf("Character %s has an expression with a space, please remove it\nExpression %s\n", s.Character, s.Expression)
		} else {
			buffer.WriteString(" ")
			buffer.WriteString(s.Expression)
		}
	}

	if s.Position != "" {
		buffer.WriteString(" at ")
		buffer.WriteString(s.Position)
	}

	return buffer.String()
}
