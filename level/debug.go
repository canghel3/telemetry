package level

type LevelDebug struct {
	levelType string
}

func NewDebugLevel() *LevelDebug {
	return &LevelDebug{levelType: "DEBUG"}
}

func (ld *LevelDebug) Type() string {
	return ld.levelType
}
