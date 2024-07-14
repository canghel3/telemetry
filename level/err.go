package level

type LevelError struct {
	levelType string
}

func NewErrorLevel() *LevelError {
	return &LevelError{levelType: "ERROR"}
}

func (le *LevelError) Type() string {
	return le.levelType
}
