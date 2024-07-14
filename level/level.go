package level

type Level interface {
	Type() string
}

type CustomLevel struct {
	levelType string
}

func NewCustomLevel(levelType string) *CustomLevel {
	return &CustomLevel{levelType: levelType}
}

func (c *CustomLevel) Type() string {
	return c.levelType
}

const (
	NoLevel    = 0
	LevelError = 1
	LevelWarn  = 2
	LevelInfo  = 3
	LevelDebug = 4
)

var LevelToText = map[uint8]string{
	LevelError: "ERROR",
	LevelWarn:  "WARN",
	LevelInfo:  "INFO",
	LevelDebug: "DEBUG",
	NoLevel:    "",
}
