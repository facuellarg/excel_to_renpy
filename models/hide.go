package models

type Hide struct {
	Text string
}

func (h Hide) Build() string {
	return "hide " + h.Text
}
