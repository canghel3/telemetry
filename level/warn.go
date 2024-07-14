package level

type LevelWarn struct {
	levelType string
}

func NewWarnLevel() *LevelWarn {
	return &LevelWarn{levelType: "WARN"}
}

func (wl *LevelWarn) Type() string {
	return wl.levelType
}
