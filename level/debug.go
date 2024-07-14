package level

type LevelDebug struct {
	levelType string
}

func Debug() *LevelDebug {
	return &LevelDebug{levelType: "DEBUG"}
}

func (ld *LevelDebug) Type() string {
	return ld.levelType
}
