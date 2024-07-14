package level

type NoLevel struct {
	levelType string
}

func NewNoLevel() *NoLevel {
	return &NoLevel{levelType: ""}
}

func (nl *NoLevel) Type() string {
	return nl.levelType
}
