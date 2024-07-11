package levels

//TODO: consider refactoring levels to an interface for easy addition of new levels

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
