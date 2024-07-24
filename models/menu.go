package models

import "bytes"

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
		if option.Label != "" {
			buffer.WriteString(":\n      jump " + option.Label)
		}
	}
	return buffer.String()
}
