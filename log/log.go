package log

import (
	"fmt"
	"telemetry/drivers"
	"telemetry/level"
	"time"
)

type Logger struct {
	buf          []byte
	level        uint8 //TODO: use an interface so that we can easily add other levels?
	outputDriver drivers.OutputDriver
	metadata     map[any]any
}

func Default() *Logger {
	return &Logger{
		outputDriver: drivers.ToStdout(),
		metadata:     nil,
		buf:          nil,
		level:        level.NoLevel,
	}
}

func (l *Logger) NoLevel() *Logger {
	l.level = level.NoLevel
	return l
}

func (l *Logger) Error() *Logger {
	l.level = level.LevelError
	return l
}

func (l *Logger) Warn() *Logger {
	l.level = level.LevelWarn
	return l
}

func (l *Logger) Info() *Logger {
	l.level = level.LevelInfo
	return l
}

func (l *Logger) Debug() *Logger {
	l.level = level.LevelDebug
	return l
}

func (l *Logger) CustomLevel() *Logger {
	l.level = 0
	return l
}

func (l *Logger) Metadata(data map[any]any) *Logger {
	l.metadata = data
	return l
}

func File(name string) *Logger {
	return &Logger{
		outputDriver: drivers.ToFileWithName(name),
	}
}

func Stdout() *Logger {
	return &Logger{
		outputDriver: drivers.ToStdout(),
	}
}

func CustomOutputDriver(driver drivers.OutputDriver) *Logger {
	return &Logger{
		outputDriver: driver,
	}
}

// Msg returns a copy of the Logger receiver, except for the buf, which is overwritten.
func (l *Logger) Msg(b []byte) *Logger {
	cpy := &Logger{
		buf:          b,
		level:        l.level,
		outputDriver: l.outputDriver,
		metadata:     l.metadata,
	}

	cpy.buf = formatLogOutput(cpy)
	return cpy
}

// Write sends the (current log buffer + received "b") to the output driver for further handling.
// The buffer is emptied and can be reused.
// If an error occurs during writing, it panics.
func (l *Logger) Write(b []byte) {
	l.buf = append(l.buf, b...)
	err := l.outputDriver.Write(formatLogOutput(l))
	if err != nil {
		panic(err)
	}

	l.buf = []byte{}
}

// TODO: enable config formatting
func formatLogOutput(l *Logger) []byte {
	//TIMESTAMP LEVEL METADATA BUFFER
	timestamp := time.Now()

	var out = make([]byte, 0)
	out = append(out, []byte(timestamp.String())...)
	out = append(out, byte(' '))
	out = append(out, []byte(level.LevelToText[l.level])...)
	out = append(out, byte(' '))

	var meta2bytes = make([]byte, 0)
	for k, v := range l.metadata {
		meta2bytes = append(meta2bytes, []byte(fmt.Sprintf("%v:%v ", k, v))...) //careful
	}
	//very careful (whitespace)
	out = append(out, meta2bytes...)
	out = append(out, l.buf...)
	out = append(out, '\n')

	return out
}
