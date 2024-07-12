package log

import (
	"fmt"
	"telemetry/drivers"
	"telemetry/levels"
	"time"
)

//TODO: add global logger so that we do not have to specify the output driver every time?

var logger = Default()

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
		level:        levels.NoLevel,
	}
}

func (l *Logger) NoLevel() *Logger {
	l.level = levels.NoLevel
	return l
}

func (l *Logger) Error() *Logger {
	l.level = levels.LevelError
	return l
}

func (l *Logger) Warn() *Logger {
	l.level = levels.LevelWarn
	return l
}

func (l *Logger) Info() *Logger {
	l.level = levels.LevelInfo
	return l
}

func (l *Logger) Debug() *Logger {
	l.level = levels.LevelDebug
	return l
}

func (l *Logger) CustomLevel() *Logger {
	l.level = 0 //TODO: implement
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

// Write sends the current log buffer to the output driver for further handling.
// The buffer is emptied and can be reused.
// If an error occurs during writing, it panics.
func (l *Logger) Write(b []byte) {
	l.buf = b
	err := l.outputDriver.Write(formatLogOutput(l))
	if err != nil {
		panic(err)
	}

	l.buf = []byte{}
}

func formatLogOutput(l *Logger) []byte {
	//TIMESTAMP LEVEL METADATA BUFFER
	timestamp := time.Now()

	var out = make([]byte, 0)
	out = append(out, []byte(timestamp.String())...)
	out = append(out, byte(' '))
	out = append(out, []byte(levels.LevelToText[l.level])...)
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
