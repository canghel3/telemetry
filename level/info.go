package level

type LevelInfo struct {
	levelType string
}

func NewInfoLevel() *LevelInfo {
	return &LevelInfo{levelType: "INFO"}
}

func (li *LevelInfo) Type() string {
	return li.levelType
}
