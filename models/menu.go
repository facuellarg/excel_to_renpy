package models

import (
	"bytes"
	"strings"
)

type Menu struct {
	Options []Options
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
		if option.Content != "" {
			buffer.WriteString(":\n      ")
			if strings.Contains(option.Content, " ") {
				buffer.WriteString("jump " + option.Content)
			}
			buffer.WriteString("\"" + option.Content + "\"")
		}
	}
	return buffer.String()
}
